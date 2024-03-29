package rabbit

import (
	"common/assert"
	"common/utils"
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"common/config"

	"sync/atomic"

	"runtime/debug"

	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

type validator interface {
	Validate() bool
}

// goodFunc verifies that the function satisfies the signature, represented as a slice of types.
// The last type is the single result type; the others are the input types.
// A final type of nil means any result type is accepted.
func goodFunc(fn reflect.Value, rtrn int, types ...reflect.Type) (e bool) {
	if fn.Kind() != reflect.Func {
		return false
	}
	// Last type is return, the rest are ins.
	if fn.Type().NumIn() != len(types)-rtrn || fn.Type().NumOut() != rtrn {
		return false
	}
	for i := 0; i < len(types)-rtrn; i++ {
		if !fn.Type().In(i).AssignableTo(types[i]) {
			return false
		}
	}

	var j int
	for i := len(types) - rtrn + 1; i < len(types); i++ {
		outType := types[i]
		if outType != nil && fn.Type().Out(j) != outType {
			panic(outType)
			return false
		}
		j++
	}

	return true
}

// RunWorker listen on a topic in Amqp
func RunWorker(
	jobPattern Job,
	function interface{},
	prefetch int) {

	in := reflect.ValueOf(jobPattern)

	fn := reflect.ValueOf(function)
	elemType := in.Type()

	var (
		t bool
		e = errors.New("test")
	)
	if !goodFunc(fn, 2, elemType, reflect.ValueOf(e).Type(), reflect.ValueOf(t).Type()) {
		logrus.Panic("function must be of type func(" + in.Type().Elem().String() + ") (bool, error)")
	}

	c, err := conn.Channel()
	assert.Nil(err)

	err = c.ExchangeDeclare(
		config.Config.AMQP.Exchange, // name
		"topic",                     // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // arguments
	)
	assert.Nil(err)

	q, err := c.QueueDeclare(jobPattern.GetQueue(), true, false, false, false, nil)
	assert.Nil(err)

	// prefetch count
	// **WARNING**
	// If ignore this, then there is a problem with rabbit. prefetch all jobs for this worker then.
	// the next worker get nothing at all!
	// **WARNING**
	err = c.Qos(prefetch, 0, false)
	assert.Nil(err)

	err = c.QueueBind(
		q.Name,                      // queue name
		jobPattern.GetTopic(),       // routing key
		config.Config.AMQP.Exchange, // exchange
		false,
		nil,
	)
	assert.Nil(err)

	consumerTag := <-utils.ID
	delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
	assert.Nil(err)

	logrus.Debug("Worker started")
	consume(delivery, jobPattern, fn, c, consumerTag)
}

func consume(
	delivery <-chan amqp.Delivery,
	jobPattern interface{},
	fn reflect.Value,
	c *amqp.Channel,
	consumerTag string) {
	waiter := sync.WaitGroup{}
	atomic.SwapInt64(&hasConsumer, 1)
bigLoop:
	for {
		select {
		case job := <-delivery:
			cp := reflect.New(reflect.TypeOf(jobPattern)).Elem().Addr().Interface()
			err := json.Unmarshal(job.Body, cp)
			if err != nil {
				logrus.Debugf("invalid job, error was : %s", err)
				assert.Nil(job.Reject(false))
				break
			}
			if v, ok := cp.(validator); ok {
				if !v.Validate() {
					_ = job.Reject(false)
					logrus.Warn("Validation failed")
				}
			}

			input := []reflect.Value{reflect.ValueOf(cp).Elem()}
			waiter.Add(1)
			go func() {
				defer waiter.Done()
				defer func() {
					if e := recover(); e != nil {
						//Panic??
						logrus.Error(e)
						debug.PrintStack()
						_ = job.Reject(false)
					}
				}()

				out := fn.Call(input)
				if out[1].Interface() == nil || out[1].Interface().(error) == nil {
					assert.Nil(job.Ack(false))
				} else {
					logrus.Debug(out[1].Interface().(error))
					assert.Nil(job.Nack(false, out[0].Interface().(bool)))
				}
			}()
		case ok := <-quit:
			_ = c.Cancel(consumerTag, false)
			waiter.Wait()
			finalizeWait()
			ok <- struct{}{}
			break bigLoop
		}

	}
}
