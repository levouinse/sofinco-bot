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

// Asah Otak Plugin
type AsahOtakPlugin struct {
	plugins.BasePlugin
	games map[int64]*AsahOtakGame
}

type AsahOtakGame struct {
	Question string
	Answer   string
	Timeout  time.Time
}

var AsahOtakInstance *AsahOtakPlugin

func init() {
	AsahOtakInstance = &AsahOtakPlugin{
		games: make(map[int64]*AsahOtakGame),
	}
	plugins.Register(AsahOtakInstance)
}

func (p *AsahOtakPlugin) Commands() []string { return []string{"asahotak"} }
func (p *AsahOtakPlugin) Tags() []string     { return []string{"game"} }
func (p *AsahOtakPlugin) Help() string       { return "Game asah otak" }
func (p *AsahOtakPlugin) RequireLimit() bool { return false }

func (p *AsahOtakPlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/asahotak?apikey=%s", ctx.Config.APIKey)
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

	game := &AsahOtakGame{
		Question: apiResponse.Result.Soal,
		Answer:   apiResponse.Result.Jawaban,
		Timeout:  time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	msg := fmt.Sprintf("🧠 *Asah Otak*\n\n%s\n\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah", game.Question)
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

func (p *AsahOtakPlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
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

// Siapakah Aku Plugin
type SiapakahAkuPlugin struct {
	plugins.BasePlugin
	games map[int64]*SiapakahAkuGame
}

type SiapakahAkuGame struct {
	Question string
	Answer   string
	Timeout  time.Time
}

var SiapakahAkuInstance *SiapakahAkuPlugin

func init() {
	SiapakahAkuInstance = &SiapakahAkuPlugin{
		games: make(map[int64]*SiapakahAkuGame),
	}
	plugins.Register(SiapakahAkuInstance)
}

func (p *SiapakahAkuPlugin) Commands() []string { return []string{"siapakahaku"} }
func (p *SiapakahAkuPlugin) Tags() []string     { return []string{"game"} }
func (p *SiapakahAkuPlugin) Help() string       { return "Game tebak siapa" }
func (p *SiapakahAkuPlugin) RequireLimit() bool { return false }

func (p *SiapakahAkuPlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/siapakahaku?apikey=%s", ctx.Config.APIKey)
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

	game := &SiapakahAkuGame{
		Question: apiResponse.Result.Soal,
		Answer:   apiResponse.Result.Jawaban,
		Timeout:  time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	msg := fmt.Sprintf("❓ *Siapakah Aku?*\n\n%s\n\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah", game.Question)
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

func (p *SiapakahAkuPlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
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

// Tebak Lagu Plugin
type TebakLaguPlugin struct {
	plugins.BasePlugin
	games map[int64]*TebakLaguGame
}

type TebakLaguGame struct {
	Audio   string
	Answer  string
	Timeout time.Time
}

var TebakLaguInstance *TebakLaguPlugin

func init() {
	TebakLaguInstance = &TebakLaguPlugin{
		games: make(map[int64]*TebakLaguGame),
	}
	plugins.Register(TebakLaguInstance)
}

func (p *TebakLaguPlugin) Commands() []string { return []string{"tebaklagu"} }
func (p *TebakLaguPlugin) Tags() []string     { return []string{"game"} }
func (p *TebakLaguPlugin) Help() string       { return "Tebak judul lagu" }
func (p *TebakLaguPlugin) RequireLimit() bool { return false }

func (p *TebakLaguPlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID

	if _, exists := p.games[chatID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Masih ada permainan yang belum selesai!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/tebaklagu?apikey=%s", ctx.Config.APIKey)
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
			Audio   string `json:"audio"`
			Jawaban string `json:"jawaban"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	game := &TebakLaguGame{
		Audio:   apiResponse.Result.Audio,
		Answer:  apiResponse.Result.Jawaban,
		Timeout: time.Now().Add(60 * time.Second),
	}
	p.games[chatID] = game

	audio := tgbotapi.NewAudio(chatID, tgbotapi.FileURL(game.Audio))
	audio.Caption = "🎵 *Tebak Lagu*\n\nApa judul lagu ini?\n⏰ Waktu: 60 detik\n💡 Ketik 'nyerah' untuk menyerah"
	audio.ParseMode = "Markdown"
	ctx.API.Send(audio)

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

func (p *TebakLaguPlugin) CheckAnswer(chatID int64, answer string) (string, bool) {
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
