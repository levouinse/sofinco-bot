package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/config"
	"github.com/levouinse/sofinco-bot/internal/database"
	"github.com/levouinse/sofinco-bot/internal/handlers"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	db       *database.Database
	config   *config.Config
	handlers *handlers.Handler
}

func New(api *tgbotapi.BotAPI, db *database.Database, cfg *config.Config) *Bot {
	b := &Bot{
		api:    api,
		db:     db,
		config: cfg,
	}
	b.handlers = handlers.New(api, db, cfg)
	return b
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			go b.handleUpdate(update)
		}
	}()
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if update.Message != nil {
		b.handlers.HandleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.handlers.HandleCallback(update.CallbackQuery)
	}
}
