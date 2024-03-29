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
	// RegisterUserHandlerWithExp like above func but runs expired func
	RegisterUserHandlerWithExp(int64, HandleMessage, func(), time.Duration)
	// UnRegisterUserHandler redirect all user message to a chat
	UnRegisterUserHandler(int64)
	// Send a message using this interface to a user
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	// GetBot return the current bot
	GetBot() *tgbotapi.BotAPI
	// NewKeyboard shows keyboard
	NewKeyboard([]string) tgbotapi.ReplyKeyboardMarkup
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
	bot      *tgbotapi.BotAPI
}

// NewTelegramBot return the telegram bot api
func NewTelegramBot(token string) TelegramBot {
	return &telegramBot{
		token:    token,
		lock:     &sync.RWMutex{},
		commands: make(map[string]HandleMessage),
		users:    make(map[int64]HandleMessage),
		sessions: make(map[int64]string),
	}
}

// TODO just one row for now
// ShowKeyboard shows keyboard to user
func (tb *telegramBot) NewKeyboard(buttonsName []string) tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{}
	for i := range buttonsName {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(buttonsName[i]))
	}
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(buttons...))
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

// RegisterUserHandlerWithExp redirect all user message to a chat
func (tb *telegramBot) RegisterUserHandlerWithExp(cid int64, hh HandleMessage, exp func(), t time.Duration) {
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
			exp()
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

func (tb *telegramBot) internalStart(bot *tgbotapi.BotAPI, updates <-chan tgbotapi.Update) {
	wg := sync.WaitGroup{}
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
						logrus.WithField("error", e).Warn(e, data)
						if config.Config.Redmine.Active {
							go utils.RedmineDoError(e, []byte(data))
						}

						if config.Config.Slack.Active {
							go utils.SlackDoMessage(e, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
						}
					}
				}()
				h(bot, update.Message)
			}()
			tb.lock.RUnlock()
			continue bigLoop
		}
		// currently we only support messages
		if !update.Message.IsCommand() {
			continue
		}
		txt := strings.Trim(strings.ToLower(update.Message.Text), " \t\n")
		for i := range tb.commands {
			if hasPrefix(txt, i) {
				wg.Add(1)
				go func(cmd string) {
					defer wg.Done()
					defer func() {
						if e := recover(); e != nil {
							stack := debug.Stack()
							dump, _ := json.MarshalIndent(update, "\t", "\t")
							data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
							logrus.WithField("error", e).Warn(e, data)
							if config.Config.Redmine.Active {
								go utils.RedmineDoError(e, []byte(data))
							}

							if config.Config.Slack.Active {
								go utils.SlackDoMessage(e, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
							}
						}
					}()
					logrus.Warn(cmd)
					tb.commands[cmd](bot, update.Message)

				}(i)
			}
		}
		tb.lock.RUnlock()
	}

	wg.Wait()
}

func (tb *telegramBot) Start() error {
	if !atomic.CompareAndSwapInt32(&tb.started, 0, 1) {
		return errors.New("already started")
	}
	defer atomic.SwapInt32(&tb.started, 0)

	var err error
	tb.bot, err = tgbotapi.NewBotAPI(tb.token)
	if err != nil {
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tb.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	tb.internalStart(tb.bot, updates)

	return nil
}

func (tb *telegramBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if atomic.CompareAndSwapInt32(&tb.started, 0, 1) {
		return tgbotapi.Message{}, errors.New("not yet started")
	}

	return tb.bot.Send(c)
}

func (tb *telegramBot) GetBot() *tgbotapi.BotAPI {
	return tb.bot
}

func hasPrefix(a, b string) bool {
	if strings.HasPrefix(a, b) {
		g := strings.Split(a, "_")
		if len(g) > 0 && g[0] != b {
			return false
		}
		return true
	}
	return false
}
