package bot

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rmarsu/easy-tg/src/waiter"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	logger *logrus.Logger
	mu     sync.Mutex
	*Router
	Waiter *waiter.Waiter[int64, tgbotapi.Update]
}

type Router struct {
	Handlers map[string]func(bot *Bot, upd *tgbotapi.Update)
}

func New(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		Waiter: waiter.New[int64, tgbotapi.Update](),
		Router: &Router{
			Handlers: make(map[string]func(bot *Bot, upd *tgbotapi.Update)),
		},
		mu:     sync.Mutex{},
		logger: logrus.New(),
		bot:    bot,
	}, nil
}

func (b *Bot) Add(call string, handler func(bot *Bot, upd *tgbotapi.Update)) {
	b.Router.Handlers[call] = handler
}

func (b *Bot) Get(call string) (func(bot *Bot, upd *tgbotapi.Update), bool) {
	if val, ok := b.Router.Handlers[call]; ok {
		return val, true
	}
	return nil, false
}

func (b *Bot) Send(c tgbotapi.Chattable) error {
	if _, err := b.bot.Send(c); err != nil {
		return err
	}
	return nil
}

func (b *Bot) WaitForMessage(upd *tgbotapi.Update) *tgbotapi.Update {
	userID := upd.Message.From.ID
	ch := b.Waiter.Add(userID)
	defer b.Waiter.Remove(userID)

	select {
	case newUpdate := <-ch:
		return &newUpdate
	case <-time.After(5 * time.Minute):
		return nil
	}
}

func (b *Bot) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		if ch := b.Waiter.Get(userID); ch != nil {
			ch <- update
			continue
		}

		handler, ok := b.Get(update.Message.Text)
		if !ok {
			b.logger.Infof("Unknown handler: %v", update.Message)
			continue
		}
		go handler(b, &update)
		b.logger.Infof("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)
	}
}
