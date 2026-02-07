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

// Tebak Anime Plugin
type TebakAnimePlugin struct {
	plugins.BasePlugin
	games map[int64]*TebakAnimeGame
}

type TebakAnimeGame struct {
	Image   string
	Answer  string
	Timeout time.Time
}

var TebakAnimeInstance *TebakAnimePlugin

func init() {
	TebakAnimeInstance = &TebakAnimePlugin{
		games: make(map[int64]*TebakAnimeGame),
	}
	plugins.Register(TebakAnimeInstance)
}

func (p *TebakAnimePlugin) Commands() []string { return []string{"tebakanime"} }
func (p *TebakAnimePlugin) Tags() []string     { return []string{"game"} }
func (p *TebakAnimePlugin) Help() string       { return "Tebak nama anime dari gambar" }
func (p *TebakAnimePlugin) RequireLimit() bool { return false }

func (p *TebakAnimePlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/tebakanime?apikey=%s", ctx.Config.APIKey)
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
			Image  string `json:"image"`
			Jawaban string `json:"jawaban"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	game := &TebakAnimeGame{
		Image:   apiResponse.Result.Image,
		Answer:  apiResponse.Result.Jawaban,
		Timeout: time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(game.Image))
	photo.Caption = "🎌 *Tebak Anime*\n\nSiapa nama anime ini?\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah"
	photo.ParseMode = "Markdown"
	ctx.API.Send(photo)

	go func() {
		time.Sleep(60 * time.Second)
		if g, exists := p.games[chatID]; exists {
			delete(p.games, chatID)
			timeout := tgbotapi.NewMessage(chatID, fmt.Sprintf("⏰ Waktu habis!\n\n*Jawaban:* %s", g.Answer))
			timeout.ParseMode = "Markdown"
			ctx.API.Send(timeout)
		}
	}()

	return nil
}

func (p *TebakAnimePlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
	game, exists := p.games[chatID]
	if !exists {
		return "", false
	}

	if time.Now().After(game.Timeout) {
		delete(p.games, chatID)
		return "⏰ Waktu habis!", true
	}

	if strings.ToLower(answer) == "nyerah" {
		delete(p.games, chatID)
		return fmt.Sprintf("🏳️ Menyerah!\n\n*Jawaban:* %s", game.Answer), true
	}

	if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(game.Answer)) {
		delete(p.games, chatID)
		return fmt.Sprintf("✅ Benar! Jawabannya adalah *%s*\n\n+50 XP", game.Answer), true
	}

	return "", false
}

// Tebak Gambar Plugin
type TebakGambarPlugin struct {
	plugins.BasePlugin
	games map[int64]*TebakGambarGame
}

type TebakGambarGame struct {
	Image   string
	Answer  string
	Timeout time.Time
}

var TebakGambarInstance *TebakGambarPlugin

func init() {
	TebakGambarInstance = &TebakGambarPlugin{
		games: make(map[int64]*TebakGambarGame),
	}
	plugins.Register(TebakGambarInstance)
}

func (p *TebakGambarPlugin) Commands() []string { return []string{"tebakgambar"} }
func (p *TebakGambarPlugin) Tags() []string     { return []string{"game"} }
func (p *TebakGambarPlugin) Help() string       { return "Tebak kata dari gambar" }
func (p *TebakGambarPlugin) RequireLimit() bool { return false }

func (p *TebakGambarPlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/tebakgambar?apikey=%s", ctx.Config.APIKey)
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
			Image  string `json:"image"`
			Jawaban string `json:"jawaban"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	game := &TebakGambarGame{
		Image:   apiResponse.Result.Image,
		Answer:  apiResponse.Result.Jawaban,
		Timeout: time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(game.Image))
	photo.Caption = "🖼️ *Tebak Gambar*\n\nApa yang ada di gambar ini?\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah"
	photo.ParseMode = "Markdown"
	ctx.API.Send(photo)

	go func() {
		time.Sleep(60 * time.Second)
		if g, exists := p.games[chatID]; exists {
			delete(p.games, chatID)
			timeout := tgbotapi.NewMessage(chatID, fmt.Sprintf("⏰ Waktu habis!\n\n*Jawaban:* %s", g.Answer))
			timeout.ParseMode = "Markdown"
			ctx.API.Send(timeout)
		}
	}()

	return nil
}

func (p *TebakGambarPlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
	game, exists := p.games[chatID]
	if !exists {
		return "", false
	}

	if time.Now().After(game.Timeout) {
		delete(p.games, chatID)
		return "⏰ Waktu habis!", true
	}

	if strings.ToLower(answer) == "nyerah" {
		delete(p.games, chatID)
		return fmt.Sprintf("🏳️ Menyerah!\n\n*Jawaban:* %s", game.Answer), true
	}

	if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(game.Answer)) {
		delete(p.games, chatID)
		return fmt.Sprintf("✅ Benar! Jawabannya adalah *%s*\n\n+50 XP", game.Answer), true
	}

	return "", false
}

// Tebak Kata Plugin
type TebakKataPlugin struct {
	plugins.BasePlugin
	games map[int64]*TebakKataGame
}

type TebakKataGame struct {
	Question string
	Answer   string
	Timeout  time.Time
}

var TebakKataInstance *TebakKataPlugin

func init() {
	TebakKataInstance = &TebakKataPlugin{
		games: make(map[int64]*TebakKataGame),
	}
	plugins.Register(TebakKataInstance)
}

func (p *TebakKataPlugin) Commands() []string { return []string{"tebakkata"} }
func (p *TebakKataPlugin) Tags() []string     { return []string{"game"} }
func (p *TebakKataPlugin) Help() string       { return "Tebak kata dari petunjuk" }
func (p *TebakKataPlugin) RequireLimit() bool { return false }

func (p *TebakKataPlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/tebakkata?apikey=%s", ctx.Config.APIKey)
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
			Soal    string `json:"soal"`
			Jawaban string `json:"jawaban"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	game := &TebakKataGame{
		Question: apiResponse.Result.Soal,
		Answer:   apiResponse.Result.Jawaban,
		Timeout:  time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	msg := fmt.Sprintf("📝 *Tebak Kata*\n\n%s\n\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah", game.Question)
	reply := tgbotapi.NewMessage(chatID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	go func() {
		time.Sleep(60 * time.Second)
		if g, exists := p.games[chatID]; exists {
			delete(p.games, chatID)
			timeout := tgbotapi.NewMessage(chatID, fmt.Sprintf("⏰ Waktu habis!\n\n*Jawaban:* %s", g.Answer))
			timeout.ParseMode = "Markdown"
			ctx.API.Send(timeout)
		}
	}()

	return nil
}

func (p *TebakKataPlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
	game, exists := p.games[chatID]
	if !exists {
		return "", false
	}

	if time.Now().After(game.Timeout) {
		delete(p.games, chatID)
		return "⏰ Waktu habis!", true
	}

	if strings.ToLower(answer) == "nyerah" {
		delete(p.games, chatID)
		return fmt.Sprintf("🏳️ Menyerah!\n\n*Jawaban:* %s", game.Answer), true
	}

	if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(game.Answer)) {
		delete(p.games, chatID)
		return fmt.Sprintf("✅ Benar! Jawabannya adalah *%s*\n\n+50 XP", game.Answer), true
	}

	return "", false
}
