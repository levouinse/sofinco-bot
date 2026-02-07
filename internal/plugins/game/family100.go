package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type Family100Plugin struct {
	plugins.BasePlugin
	games map[int64]*Family100Game
}

type Family100Game struct {
	Question string
	Answers  []string
	Answered []bool
	Timeout  time.Time
}

var Family100Instance *Family100Plugin

func init() {
	Family100Instance = &Family100Plugin{
		games: make(map[int64]*Family100Game),
	}
	plugins.Register(Family100Instance)
}

func (p *Family100Plugin) Commands() []string { return []string{"family100"} }
func (p *Family100Plugin) Tags() []string     { return []string{"game"} }
func (p *Family100Plugin) Help() string       { return "Family 100 game" }
func (p *Family100Plugin) RequireLimit() bool { return false }

func (p *Family100Plugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/family100?apikey=%s", ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan soal: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Soal    string   `json:"soal"`
			Jawaban []string `json:"jawaban"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w. Response: %s", err, string(body))
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	if apiResponse.Result.Soal == "" || len(apiResponse.Result.Jawaban) == 0 {
		return fmt.Errorf("soal kosong dari API")
	}

	game := &Family100Game{
		Question: apiResponse.Result.Soal,
		Answers:  apiResponse.Result.Jawaban,
		Answered: make([]bool, len(apiResponse.Result.Jawaban)),
		Timeout:  time.Now().Add(3 * time.Minute),
	}
	p.games[chatID] = game

	msg := fmt.Sprintf("🎯 *Family 100*\n\n*Soal:* %s\n\nTerdapat *%d* jawaban\n⏰ Waktu: 3 menit\n💡 Ketik 'nyerah' untuk menyerah",
		game.Question, len(game.Answers))

	reply := tgbotapi.NewMessage(chatID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	go func() {
		time.Sleep(3 * time.Minute)
		if g, exists := p.games[chatID]; exists {
			delete(p.games, chatID)
			allAnswers := ""
			for i, ans := range g.Answers {
				allAnswers += fmt.Sprintf("%d. %s\n", i+1, ans)
			}
			timeout := tgbotapi.NewMessage(chatID, fmt.Sprintf("⏰ Waktu habis!\n\n*Jawaban yang benar:*\n%s", allAnswers))
			timeout.ParseMode = "Markdown"
			ctx.API.Send(timeout)
		}
	}()

	return nil
}

func (p *Family100Plugin) CheckAnswer(chatID int64, userID int64, answer string) (string, bool) {
	game, exists := p.games[chatID]
	if !exists {
		return "", false
	}

	if time.Now().After(game.Timeout) {
		delete(p.games, chatID)
		return "⏰ Waktu habis!", true
	}

	if strings.ToLower(answer) == "nyerah" {
		allAnswers := ""
		for i, ans := range game.Answers {
			allAnswers += fmt.Sprintf("%d. %s\n", i+1, ans)
		}
		delete(p.games, chatID)
		return fmt.Sprintf("🏳️ Menyerah!\n\n*Jawaban yang benar:*\n%s", allAnswers), true
	}

	// Check if answer is correct
	for i, ans := range game.Answers {
		if !game.Answered[i] && strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(ans)) {
			game.Answered[i] = true

			// Count remaining answers
			remaining := 0
			for _, answered := range game.Answered {
				if !answered {
					remaining++
				}
			}

			if remaining == 0 {
				delete(p.games, chatID)
				return fmt.Sprintf("✅ Benar! *%s*\n\n🎉 Semua jawaban telah ditemukan!", ans), true
			}

			return fmt.Sprintf("✅ Benar! *%s*\n\n💡 Masih ada %d jawaban lagi", ans, remaining), true
		}
	}

	return "", false
}
