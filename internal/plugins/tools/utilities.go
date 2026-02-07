package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

// Sticker Maker Plugin
type StickerPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &StickerPlugin{}
	plugins.Register(p)
}

func (p *StickerPlugin) Commands() []string { return []string{"sticker", "s"} }
func (p *StickerPlugin) Tags() []string     { return []string{"tools"} }
func (p *StickerPlugin) Help() string       { return "Convert image/video to sticker" }
func (p *StickerPlugin) RequireLimit() bool { return true }

func (p *StickerPlugin) Execute(ctx *plugins.Context) error {
	if ctx.Message.ReplyToMessage == nil {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Reply to an image or video to convert to sticker!"))
		return nil
	}

	var fileID string
	if ctx.Message.ReplyToMessage.Photo != nil && len(ctx.Message.ReplyToMessage.Photo) > 0 {
		fileID = ctx.Message.ReplyToMessage.Photo[len(ctx.Message.ReplyToMessage.Photo)-1].FileID
	} else if ctx.Message.ReplyToMessage.Video != nil {
		fileID = ctx.Message.ReplyToMessage.Video.FileID
	} else {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "❌ Reply to an image or video!"))
		return nil
	}

	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Converting to sticker...")
	sent, _ := ctx.API.Send(msg)

	// Get file
	file, err := ctx.API.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to get file")
		ctx.API.Send(edit)
		return err
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", ctx.Config.BotToken, file.FilePath)

	// Send as sticker
	sticker := tgbotapi.NewSticker(ctx.Message.Chat.ID, tgbotapi.FileURL(fileURL))
	if _, err := ctx.API.Send(sticker); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to create sticker")
		ctx.API.Send(edit)
		return err
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
	return nil
}

// ToImage Plugin
type ToImagePlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &ToImagePlugin{}
	plugins.Register(p)
}

func (p *ToImagePlugin) Commands() []string { return []string{"toimg", "toimage"} }
func (p *ToImagePlugin) Tags() []string     { return []string{"tools"} }
func (p *ToImagePlugin) Help() string       { return "Convert sticker to image" }
func (p *ToImagePlugin) RequireLimit() bool { return true }

func (p *ToImagePlugin) Execute(ctx *plugins.Context) error {
	if ctx.Message.ReplyToMessage == nil || ctx.Message.ReplyToMessage.Sticker == nil {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Reply to a sticker to convert to image!"))
		return nil
	}

	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Converting to image...")
	sent, _ := ctx.API.Send(msg)

	fileID := ctx.Message.ReplyToMessage.Sticker.FileID
	file, err := ctx.API.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to get file")
		ctx.API.Send(edit)
		return err
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", ctx.Config.BotToken, file.FilePath)

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(fileURL))
	if _, err := ctx.API.Send(photo); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to convert")
		ctx.API.Send(edit)
		return err
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
	return nil
}

// QR Code Generator
type QRCodePlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &QRCodePlugin{}
	plugins.Register(p)
}

func (p *QRCodePlugin) Commands() []string { return []string{"qrcode", "qr"} }
func (p *QRCodePlugin) Tags() []string     { return []string{"tools"} }
func (p *QRCodePlugin) Help() string       { return "Generate QR code" }
func (p *QRCodePlugin) RequireLimit() bool { return false }

func (p *QRCodePlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /qrcode <text>"))
		return nil
	}

	text := ctx.Message.Text[len("/qrcode "):]
	qrURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=500x500&data=%s", text)

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(qrURL))
	photo.Caption = fmt.Sprintf("📱 QR Code for:\n%s", text)
	ctx.API.Send(photo)

	return nil
}

// Translate Plugin
type TranslatePlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &TranslatePlugin{}
	plugins.Register(p)
}

func (p *TranslatePlugin) Commands() []string { return []string{"translate", "tr"} }
func (p *TranslatePlugin) Tags() []string     { return []string{"tools"} }
func (p *TranslatePlugin) Help() string       { return "Translate text" }
func (p *TranslatePlugin) RequireLimit() bool { return false }

func (p *TranslatePlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) < 2 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /translate <lang> <text>\n\nExample: /translate en Halo dunia"))
		return nil
	}

	lang := ctx.Args[0]
	text := ctx.Message.Text[len("/translate ")+len(lang)+1:]

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/tools/translate?text=%s&lang=%s&apikey=%s", text, lang, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("gagal translate: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse response: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	msg := fmt.Sprintf("🌐 *Translation*\n\n*Original:*\n%s\n\n*Translated (%s):*\n%s", text, lang, apiResponse.Result)
	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	return nil
}

// Wikipedia Plugin
type WikipediaPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &WikipediaPlugin{}
	plugins.Register(p)
}

func (p *WikipediaPlugin) Commands() []string { return []string{"wikipedia", "wiki"} }
func (p *WikipediaPlugin) Tags() []string     { return []string{"tools"} }
func (p *WikipediaPlugin) Help() string       { return "Search Wikipedia" }
func (p *WikipediaPlugin) RequireLimit() bool { return false }

func (p *WikipediaPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /wikipedia <query>"))
		return nil
	}

	query := ctx.Message.Text[len("/wikipedia "):]
	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/search/wikipedia?query=%s&apikey=%s", query, ctx.Config.APIKey)
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("gagal search: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Title   string `json:"title"`
			Extract string `json:"extract"`
			URL     string `json:"url"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse response: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	msg := fmt.Sprintf("📚 *%s*\n\n%s\n\n🔗 [Read more](%s)", 
		apiResponse.Result.Title, 
		apiResponse.Result.Extract, 
		apiResponse.Result.URL)
	
	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	reply.DisableWebPagePreview = true
	ctx.API.Send(reply)

	return nil
}

// Calculator Plugin
type CalculatorPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &CalculatorPlugin{}
	plugins.Register(p)
}

func (p *CalculatorPlugin) Commands() []string { return []string{"calc", "calculator"} }
func (p *CalculatorPlugin) Tags() []string     { return []string{"tools"} }
func (p *CalculatorPlugin) Help() string       { return "Calculate math expression" }
func (p *CalculatorPlugin) RequireLimit() bool { return false }

func (p *CalculatorPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /calc <expression>\n\nExample: /calc 2+2*3"))
		return nil
	}

	expression := ctx.Message.Text[len("/calc "):]
	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/tools/calculator?q=%s&apikey=%s", expression, ctx.Config.APIKey)
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("gagal calculate: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("gagal parse response: %w", err)
	}

	if !apiResponse.Status {
		return fmt.Errorf("API error: %s", apiResponse.Message)
	}

	msg := fmt.Sprintf("🧮 *Calculator*\n\n`%s` = `%s`", expression, apiResponse.Result)
	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	return nil
}
