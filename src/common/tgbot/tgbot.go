package tgbot

import (
	"common/config"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// TelegramBot is an interface to handle the telegram bot
type TelegramBot interface {
	// RegisterHandler try to register a handler in system, the first is the command to match the
	// next arg is the handler function
	RegisterMessageHandler(string, HandleMessage) error
	// Start the handler
	Start() error
	// RegisterUserHandler redirect all user message to a chat
	RegisterUserHandler(int64, HandleMessage, time.Duration)
	// UnRegisterUserHandler redirect all user message to a chat
	UnRegisterUserHandler(int64)
}

// HandleMessage is the callback for message with router
type HandleMessage func(*tgbotapi.BotAPI, *tgbotapi.Message)

type telegramBot struct {
	token    string
	lock     *sync.RWMutex
	commands map[string]HandleMessage
	started  int32
	users    map[int64]HandleMessage
	sessions map[int64]string
}

// NewTelegramBot return the telegram bot api
func NewTelegramBot(token string) TelegramBot {
	return &telegramBot{
		token:    token,
		lock:     &sync.RWMutex{},
		commands: make(map[string]HandleMessage),
		users:    make(map[int64]HandleMessage),
	}
}

func (tb *telegramBot) RegisterMessageHandler(cmd string, handler HandleMessage) error {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	cmd = strings.Trim(strings.ToLower(cmd), " \t\n")
	if _, ok := tb.commands[cmd]; ok {
		return fmt.Errorf("already registered : %s", cmd)
	}
	tb.commands[cmd] = handler
	return nil
}

// RegisterUserHandler redirect all user message to a chat
func (tb *telegramBot) RegisterUserHandler(cid int64, hh HandleMessage, t time.Duration) {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	tb.users[cid] = hh
	lockID := <-utils.ID
	tb.sessions[cid] = lockID

	go func() {
		<-time.After(t)
		tb.lock.Lock()
		defer tb.lock.Unlock()

		if tb.sessions[cid] == lockID {
			delete(tb.sessions, cid)
			delete(tb.users, cid)
		}
	}()
}

// UnRegisterUserHandler redirect all user message to a chat
func (tb *telegramBot) UnRegisterUserHandler(cid int64) {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	delete(tb.sessions, cid)
	delete(tb.users, cid)
}

func (tb *telegramBot) Start() error {
	if !atomic.CompareAndSwapInt32(&tb.started, 0, 1) {
		return errors.New("already started")
	}
	defer atomic.SwapInt32(&tb.started, 0)

	wg := sync.WaitGroup{}
	bot, err := tgbotapi.NewBotAPI(tb.token)
	if err != nil {
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

bigLoop:
	for update := range updates {
		// currently we only support messages
		if update.Message == nil {
			continue
		}
		tb.lock.RLock()
		if h, ok := tb.users[update.Message.Chat.ID]; ok {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if e := recover(); e != nil {
						stack := debug.Stack()
						dump, _ := json.MarshalIndent(update, "\t", "\t")
						data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
						logrus.WithField("error", err).Warn(err, data)
						if config.Config.Redmine.Active {
							go utils.RedmineDoError(err, []byte(data))
						}

						if config.Config.Slack.Active {
							go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
						}
					}
				}()
				h(bot, update.Message)
			}()
			continue bigLoop
		}
		// currently we only support messages
		if !update.Message.IsCommand() {
			continue
		}
		txt := strings.Trim(strings.ToLower(update.Message.Text), " \t\n")
		for i := range tb.commands {
			if strings.HasPrefix(txt, i) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() {
						if e := recover(); e != nil {
							stack := debug.Stack()
							dump, _ := json.MarshalIndent(update, "\t", "\t")
							data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
							logrus.WithField("error", err).Warn(err, data)
							if config.Config.Redmine.Active {
								go utils.RedmineDoError(err, []byte(data))
							}

							if config.Config.Slack.Active {
								go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
							}
						}
					}()
					tb.commands[i](bot, update.Message)
				}()
			}
		}
		tb.lock.RUnlock()
	}

	wg.Wait()
	return nil
}
