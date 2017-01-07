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
}

// HandleMessage is the callback for message with router
type HandleMessage func(*tgbotapi.BotAPI, *tgbotapi.Message)

type telegramBot struct {
	token    string
	lock     *sync.RWMutex
	commands map[string]HandleMessage
	started  int32
}

// NewTelegramBot return the telegram bot api
func NewTelegramBot(token string) TelegramBot {
	return &telegramBot{
		token:    token,
		lock:     &sync.RWMutex{},
		commands: make(map[string]HandleMessage),
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

func (tg *telegramBot) Start() error {
	if !atomic.CompareAndSwapInt32(&tg.started, 0, 1) {
		return errors.New("already started")
	}
	defer atomic.SwapInt32(&tg.started, 0)

	wg := sync.WaitGroup{}
	bot, err := tgbotapi.NewBotAPI(tg.token)
	if err != nil {
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		// currently we only support messages
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		txt := strings.Trim(strings.ToLower(update.Message.Text), " \t\n")
		tg.lock.RLock()
		for i := range tg.commands {
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
					tg.commands[i](bot, update.Message)
				}()
			}
		}
		tg.lock.RUnlock()
	}

	wg.Wait()
	return nil
}
