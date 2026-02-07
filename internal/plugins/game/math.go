package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type MathPlugin struct {
	plugins.BasePlugin
	sessions map[int64]*MathSession
}

type MathSession struct {
	Question string
	Answer   int
	Mode     string
	Bonus    int
	Timeout  time.Time
}

var MathInstance *MathPlugin

var modes = map[string]struct {
	Bonus int
	Time  int
	Money int
}{
	"noob":        {10, 20, 500},
	"easy":        {20, 30, 1000},
	"medium":      {40, 40, 2500},
	"hard":        {100, 60, 5000},
	"master":      {250, 70, 10000},
	"grandmaster": {500, 90, 25000},
	"legendary":   {1000, 120, 50000},
	"mythic":      {3000, 150, 75000},
	"god":         {5000, 200, 100000},
}

func init() {
	MathInstance = &MathPlugin{
		sessions: make(map[int64]*MathSession),
	}
	plugins.Register(MathInstance)
}

func (p *MathPlugin) Commands() []string { return []string{"math"} }
func (p *MathPlugin) Tags() []string     { return []string{"game"} }
func (p *MathPlugin) Help() string       { return "Math quiz game" }
func (p *MathPlugin) RequireLimit() bool { return false }

func (p *MathPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		modeList := []string{}
		for k := range modes {
			modeList = append(modeList, k)
		}
		msg := fmt.Sprintf("Silakan pilih tingkat kesulitan.\nMode: %s\n\nContoh: /math medium",
			strings.Join(modeList, " | "))
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, msg))
		return nil
	}

	mode := strings.ToLower(ctx.Args[0])
	modeData, ok := modes[mode]
	if !ok {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Mode tidak ditemukan!"))
		return nil
	}

	if _, exists := p.sessions[ctx.Message.Chat.ID]; exists {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masih ada soal yang belum terjawab!"))
		return nil
	}

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/game/math?apikey=%s&mode=%s", ctx.Config.APIKey, mode)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan soal: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Parse API response
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Soal    string `json:"soal"`
			Jawaban string `json:"jawaban"`
			Level   string `json:"level"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse soal: %w. Response: %s", err, string(body))
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	if apiResponse.Result.Soal == "" {
		return fmt.Errorf("soal kosong dari API")
	}

	// Parse answer from string to int
	answer := 0
	fmt.Sscanf(apiResponse.Result.Jawaban, "%d", &answer)

	p.sessions[ctx.Message.Chat.ID] = &MathSession{
		Question: apiResponse.Result.Soal,
		Answer:   answer,
		Mode:     mode,
		Bonus:    modeData.Bonus,
		Timeout:  time.Now().Add(time.Duration(modeData.Time) * time.Millisecond),
	}

	msg := fmt.Sprintf("🧮 *Math Quiz - %s*\n\n%s = ?\n\nBonus: +%d XP\nWaktu: %d detik",
		strings.ToUpper(mode), apiResponse.Result.Soal, modeData.Bonus, modeData.Time/1000)
	
	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	go func() {
		time.Sleep(time.Duration(modeData.Time) * time.Millisecond)
		if session, exists := p.sessions[ctx.Message.Chat.ID]; exists {
			delete(p.sessions, ctx.Message.Chat.ID)
			timeout := tgbotapi.NewMessage(ctx.Message.Chat.ID, 
				fmt.Sprintf("⏰ Waktu habis! Jawaban: %d", session.Answer))
			ctx.API.Send(timeout)
		}
	}()

	return nil
}

func (p *MathPlugin) CheckAnswer(ctx *plugins.Context, answer string) bool {
	session, exists := p.sessions[ctx.Message.Chat.ID]
	if !exists {
		return false
	}

	if time.Now().After(session.Timeout) {
		delete(p.sessions, ctx.Message.Chat.ID)
		return false
	}

	userAnswer, err := strconv.Atoi(answer)
	if err != nil {
		return false
	}

	if userAnswer == session.Answer {
		delete(p.sessions, ctx.Message.Chat.ID)
		ctx.User.Exp += session.Bonus
		ctx.DB.SaveUser(ctx.User)
		
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, 
			fmt.Sprintf("✅ Benar! +%d XP", session.Bonus))
		ctx.API.Send(msg)
		return true
	}

	return false
}
