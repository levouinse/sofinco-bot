package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type SuitPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &SuitPlugin{}
	plugins.Register(p)
}

func (p *SuitPlugin) Commands() []string { return []string{"suit"} }
func (p *SuitPlugin) Tags() []string     { return []string{"game"} }
func (p *SuitPlugin) Help() string       { return "Rock Paper Scissors game" }
func (p *SuitPlugin) RequireLimit() bool { return false }

func (p *SuitPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		msg := "🎮 *Rock Paper Scissors*\n\nPilihan yang tersedia:\n• batu\n• gunting\n• kertas\n\nContoh: /suit batu"
		reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
		reply.ParseMode = "Markdown"
		ctx.API.Send(reply)
		return nil
	}

	userChoice := strings.ToLower(ctx.Args[0])
	validChoices := []string{"batu", "gunting", "kertas"}

	// Validate user choice
	valid := false
	for _, choice := range validChoices {
		if userChoice == choice {
			valid = true
			break
		}
	}

	if !valid {
		msg := "❌ Pilihan tidak valid!\n\nPilihan yang tersedia:\n• batu\n• gunting\n• kertas"
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, msg))
		return nil
	}

	// Bot random choice
	rand.Seed(time.Now().UnixNano())
	botChoice := validChoices[rand.Intn(len(validChoices))]

	// Determine winner
	result := ""
	emoji := map[string]string{
		"batu":    "🪨",
		"gunting": "✂️",
		"kertas":  "📄",
	}

	if userChoice == botChoice {
		result = "🤝 *Seri!*"
	} else if (userChoice == "batu" && botChoice == "gunting") ||
		(userChoice == "gunting" && botChoice == "kertas") ||
		(userChoice == "kertas" && botChoice == "batu") {
		result = "🎉 *Kamu Menang!*\n+1000 Money"
		ctx.User.Exp += 10
		ctx.DB.SaveUser(ctx.User)
	} else {
		result = "😔 *Kamu Kalah!*"
	}

	msg := fmt.Sprintf("🎮 *Rock Paper Scissors*\n\n%s\n\nKamu: %s %s\nBot: %s %s",
		result, emoji[userChoice], userChoice, emoji[botChoice], botChoice)

	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	return nil
}
