package handlers

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/config"
	"github.com/levouinse/sofinco-bot/internal/database"
	"github.com/levouinse/sofinco-bot/internal/downloader"
	"github.com/levouinse/sofinco-bot/internal/plugins"
	"github.com/levouinse/sofinco-bot/internal/plugins/game"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/ai"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/downloader"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/maker"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/owner"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/stalker"
	_ "github.com/levouinse/sofinco-bot/internal/plugins/tools"
)

type Handler struct {
	api        *tgbotapi.BotAPI
	db         *database.Database
	config     *config.Config
	downloader *downloader.YouTubeDownloader
	startTime  time.Time
}

func New(api *tgbotapi.BotAPI, db *database.Database, cfg *config.Config) *Handler {
	return &Handler{
		api:        api,
		db:         db,
		config:     cfg,
		downloader: downloader.NewYouTubeDownloader(cfg.APIKey, cfg.BaseAPIURL),
		startTime:  time.Now(),
	}
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) {
	if msg.Text == "" {
		return
	}

	user, err := h.db.GetOrCreateUser(msg.From.ID, msg.From.UserName, msg.From.FirstName)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	// Check for game responses first (before command parsing)
	text := strings.TrimSpace(msg.Text)
	
	// Check TicTacToe move
	if game.TicTacToeInstance != nil {
		if response, handled := game.TicTacToeInstance.HandleMove(msg.Chat.ID, msg.From.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "HTML"
			h.api.Send(reply)
			return
		}
	}

	// Check Family100 answer
	if game.Family100Instance != nil {
		if response, handled := game.Family100Instance.CheckAnswer(msg.Chat.ID, msg.From.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Math answer
	if game.MathInstance != nil {
		if game.MathInstance.CheckAnswer(&plugins.Context{
			API:     h.api,
			DB:      h.db,
			Config:  h.config,
			Message: msg,
			User:    user,
		}, text) {
			return
		}
	}

	// Check Tebak Anime answer
	if game.TebakAnimeInstance != nil {
		if response, handled := game.TebakAnimeInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Tebak Gambar answer
	if game.TebakGambarInstance != nil {
		if response, handled := game.TebakGambarInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Tebak Kata answer
	if game.TebakKataInstance != nil {
		if response, handled := game.TebakKataInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Asah Otak answer
	if game.AsahOtakInstance != nil {
		if response, handled := game.AsahOtakInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Siapakah Aku answer
	if game.SiapakahAkuInstance != nil {
		if response, handled := game.SiapakahAkuInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	// Check Tebak Lagu answer
	if game.TebakLaguInstance != nil {
		if response, handled := game.TebakLaguInstance.CheckAnswer(msg.Chat.ID, text); handled {
			reply := tgbotapi.NewMessage(msg.Chat.ID, response)
			reply.ParseMode = "Markdown"
			h.api.Send(reply)
			return
		}
	}

	command := strings.ToLower(text)
	if strings.HasPrefix(command, "/") {
		command = strings.TrimPrefix(command, "/")
		parts := strings.Fields(command)
		if len(parts) == 0 {
			return
		}
		cmd := parts[0]
		args := parts[1:]

		// Check plugin registry first
		if plugin, exists := plugins.Registry[cmd]; exists {
			ctx := &plugins.Context{
				API:     h.api,
				DB:      h.db,
				Config:  h.config,
				Message: msg,
				User:    user,
				Args:    args,
				Command: cmd,
			}

			// Check requirements
			if plugin.RequireLimit() && user.Limit <= 0 && !user.Premium {
				h.sendMessage(msg.Chat.ID, "❌ Limit Anda habis! Upgrade ke premium untuk akses unlimited.")
				return
			}

			if err := plugin.Execute(ctx); err != nil {
				h.sendMessage(msg.Chat.ID, fmt.Sprintf("❌ Error: %v", err))
			}
			return
		}

		// Built-in commands
		switch cmd {
		case "start", "menu":
			h.handleStart(msg, user)
		case "ping":
			h.handlePing(msg)
		case "getid":
			h.handleGetID(msg)
		case "play":
			h.handlePlay(msg, user)
		case "limit":
			h.handleLimit(msg, user)
		case "profile":
			h.handleProfile(msg, user)
		default:
			h.sendMessage(msg.Chat.ID, "Command not found. Use /menu to see available commands.")
		}
	}
}

func (h *Handler) handleStart(msg *tgbotapi.Message, user *database.User) {
	text := fmt.Sprintf("Hi %s! 👋\n\n"+
		"*Sofinco Bot* - Your Telegram Assistant\n\n"+
		"📊 Your Stats:\n"+
		"├ Limit: %d\n"+
		"├ XP: %d\n"+
		"└ Premium: %v\n\n"+
		"Select a category below:",
		user.FirstName, user.Limit, user.Exp, user.Premium)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📥 Downloader", "cat_downloader"),
			tgbotapi.NewInlineKeyboardButtonData("🎮 Games", "cat_games"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛠 Tools", "cat_tools"),
			tgbotapi.NewInlineKeyboardButtonData("🤖 AI", "cat_ai"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ Info", "cat_info"),
			tgbotapi.NewInlineKeyboardButtonData("👑 Owner", "cat_owner"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Statistics", "stats"),
		),
	)

	msgConfig := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgConfig.ReplyMarkup = keyboard
	msgConfig.ParseMode = "Markdown"
	h.api.Send(msgConfig)
}

func (h *Handler) handlePing(msg *tgbotapi.Message) {
	start := time.Now()
	sent := h.sendMessage(msg.Chat.ID, "Calculating...")
	elapsed := time.Since(start)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(h.startTime).Round(time.Second)

	text := fmt.Sprintf("Pong\n\nResponse Time: %dms\nMemory: %.2f MB\nUptime: %s",
		elapsed.Milliseconds(),
		float64(m.Alloc)/1024/1024,
		uptime)

	edit := tgbotapi.NewEditMessageText(msg.Chat.ID, sent.MessageID, text)
	h.api.Send(edit)
}

func (h *Handler) handleGetID(msg *tgbotapi.Message) {
	text := fmt.Sprintf("Your ID: %d\nYour Username: @%s\nChat ID: %d",
		msg.From.ID, msg.From.UserName, msg.Chat.ID)
	h.sendMessage(msg.Chat.ID, text)
}

func (h *Handler) handlePlay(msg *tgbotapi.Message, user *database.User) {
	args := strings.TrimPrefix(msg.Text, "/play ")
	if args == msg.Text || args == "" {
		h.sendMessage(msg.Chat.ID, "Usage: /play <song name or URL>")
		return
	}

	if user.Limit <= 0 && !user.Premium {
		h.sendMessage(msg.Chat.ID, "Your limit has been exhausted. Upgrade to premium for unlimited access.")
		return
	}

	waitMsg := h.sendMessage(msg.Chat.ID, "Searching and downloading...")

	result, err := h.downloader.Download(args)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(msg.Chat.ID, waitMsg.MessageID, fmt.Sprintf("Error: %v", err))
		h.api.Send(edit)
		return
	}

	h.api.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, waitMsg.MessageID))

	caption := fmt.Sprintf("Title: %s\nDuration: %s\nViews: %s\nAuthor: %s",
		result.Result.Title,
		result.Result.Duration,
		result.Result.Views,
		result.Result.Author,
	)

	photo := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileURL(result.Result.Thumbnail))
	photo.Caption = caption
	h.api.Send(photo)

	audio := tgbotapi.NewAudio(msg.Chat.ID, tgbotapi.FileURL(result.Result.MP3))
	audio.Title = result.Result.Title
	audio.Performer = result.Result.Author
	h.api.Send(audio)

	if !user.Premium {
		user.Limit--
		h.db.SaveUser(user)
		h.sendMessage(msg.Chat.ID, fmt.Sprintf("Limit used: 1\nRemaining: %d", user.Limit))
	}
}

func (h *Handler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	h.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	data := callback.Data

	switch {
	case data == "stats":
		h.handleStats(callback)
	case strings.HasPrefix(data, "cat_"):
		category := strings.TrimPrefix(data, "cat_")
		h.handleCategory(callback, category)
	case data == "back_menu":
		h.handleBackToMenu(callback)
	}
}

func (h *Handler) handleStats(callback *tgbotapi.CallbackQuery) {
	text := "Statistics\n\nUsers: 0\nChats: 0\nCommands: 0"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "back_menu"),
		),
	)

	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
	edit.ReplyMarkup = &keyboard
	h.api.Send(edit)
}

func (h *Handler) handleCategory(callback *tgbotapi.CallbackQuery, category string) {
	var text string
	var keyboard tgbotapi.InlineKeyboardMarkup

	switch category {
	case "downloader":
		text = "📥 *Downloader Commands*\n\n" +
			"/play - Download audio from YouTube\n" +
			"/ytv - Download video from YouTube\n" +
			"/tiktok - Download TikTok video\n" +
			"/ytmp3 - YouTube to MP3"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	case "games":
		text = "🎮 *Game Commands*\n\n" +
			"/math - Math quiz game\n" +
			"/tictactoe - Tic Tac Toe game\n" +
			"/suit - Rock Paper Scissors\n" +
			"/family100 - Family 100 game"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	case "tools":
		text = "🛠 *Tools Commands*\n\n" +
			"/remini - Enhance image quality\n" +
			"/getid - Get your user ID\n" +
			"/sticker - Create sticker\n" +
			"/toimg - Convert sticker to image\n\n" +
			"*Stalker:*\n" +
			"/igstalk - Instagram stalker\n" +
			"/ghstalk - GitHub stalker\n\n" +
			"*Maker:*\n" +
			"/jadianime - Convert to anime style"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	case "ai":
		text = "🤖 *AI Commands*\n\n" +
			"/ai - Chat with AI\n" +
			"/ask - Ask AI anything\n" +
			"/chatgpt - ChatGPT"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	case "info":
		text = "ℹ️ *Info Commands*\n\n" +
			"/ping - Bot status\n" +
			"/stats - Bot statistics\n" +
			"/limit - Check your limit\n" +
			"/profile - Your profile"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	case "owner":
		text = "👑 *Owner Commands*\n\n" +
			"/broadcast - Broadcast message\n" +
			"/addprem - Add premium user\n" +
			"/stats - Bot statistics\n\n" +
			"Restricted to bot owner only"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	default:
		text = "Category not found"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back", "back_menu"),
			),
		)
	}

	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
	edit.ReplyMarkup = &keyboard
	edit.ParseMode = "Markdown"
	h.api.Send(edit)
}

func (h *Handler) handleBackToMenu(callback *tgbotapi.CallbackQuery) {
	text := "*Sofinco Bot* - Your Telegram Assistant\n\nSelect a category below:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📥 Downloader", "cat_downloader"),
			tgbotapi.NewInlineKeyboardButtonData("🎮 Games", "cat_games"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛠 Tools", "cat_tools"),
			tgbotapi.NewInlineKeyboardButtonData("🤖 AI", "cat_ai"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ Info", "cat_info"),
			tgbotapi.NewInlineKeyboardButtonData("👑 Owner", "cat_owner"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Statistics", "stats"),
		),
	)

	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
	edit.ReplyMarkup = &keyboard
	edit.ParseMode = "Markdown"
	h.api.Send(edit)
}

func (h *Handler) sendMessage(chatID int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	sent, err := h.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	return sent
}

func (h *Handler) handleLimit(msg *tgbotapi.Message, user *database.User) {
	text := fmt.Sprintf("📊 *Limit Information*\n\n"+
		"Your Limit: %d\n"+
		"Premium: %v\n\n"+
		"Limit resets daily at 00:00 WIB\n"+
		"Upgrade to premium for unlimited access!",
		user.Limit, user.Premium)
	
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	h.api.Send(reply)
}

func (h *Handler) handleProfile(msg *tgbotapi.Message, user *database.User) {
	text := fmt.Sprintf("👤 *Your Profile*\n\n"+
		"Name: %s\n"+
		"Username: @%s\n"+
		"ID: %d\n\n"+
		"📊 Stats:\n"+
		"├ XP: %d\n"+
		"├ Level: %d\n"+
		"├ Limit: %d\n"+
		"└ Premium: %v\n\n"+
		"Registered: %v",
		user.FirstName, user.Username, user.ID,
		user.Exp, user.Level, user.Limit, user.Premium,
		user.RegisteredAt.Format("2006-01-02"))
	
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	h.api.Send(reply)
}
