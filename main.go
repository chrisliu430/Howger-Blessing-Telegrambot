package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	TelegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// --- Variable
	var err error
	var botToken string
	var format string
	var bot *TelegramBotAPI.BotAPI
	var updates TelegramBotAPI.UpdatesChannel
	// --- Bot API Setting
	port := os.Getenv("PORT")
	url := os.Getenv("URL")
	botToken = os.Getenv("Token")
	addr := fmt.Sprintf(":%s", port)
	format = "mp4"
	bot, err = TelegramBotAPI.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	// --- Handle Webhook Function
	_, err = bot.SetWebhook(TelegramBotAPI.NewWebhookWithCert(url, nil))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
	}
	updates = bot.ListenForWebhook("/")
	// ---
	go http.ListenAndServe(addr, nil)
	botProcessMessage(bot, updates, format)
}

func botReponseMessage(bot *TelegramBotAPI.BotAPI, text string, chatID int64, format string) {
	msg := TelegramBotAPI.NewMessage(chatID, "")
	result := cralwerToGetVideo(text, format)
	log.Println(result)
	if result["error"] != "" {
		msg.Text = "Howger沒念這個字哦"
		bot.Send(msg)
	} else if format == "mp4" {
		videoURL := "http://howfun.macs1207.info/video?v=" + result["media_id"]
		videoMsg := TelegramBotAPI.NewVideoShare(chatID, videoURL)
		bot.Send(videoMsg)
	} else if format == "mp3" {
		audioURL := "http://howfun.macs1207.info/audio?a=" + result["media_id"]
		audioMsg := TelegramBotAPI.NewAudioShare(chatID, audioURL)
		bot.Send(audioMsg)
	}
}

func botProcessMessage(bot *TelegramBotAPI.BotAPI, updates TelegramBotAPI.UpdatesChannel, format string) {
	for update := range updates {
		callGoogleAnalytics()
		msg := TelegramBotAPI.NewMessage(update.Message.Chat.ID, "")
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case "start":
			msg.Text = "直接輸入文字即會傳送影片給你\n若僅需要聲音請使用/voice [文字]來取得\n/help 可以查看命令\nHow哥並無唸英文，所以可以打相似的音來讓HOW哥念:)"
		case "help":
			msg.Text = "--- Help List ---\n/start 可以查看如何使用\n/feedback 可以填寫建議(未完成)\n/audio 變更檔案格式為音訊檔\n/video 變更檔案格式為影片檔\n/info 提供作者群、作者聯繫信箱及開源連結"
		case "info":
			msg.Text = "作者: C.H(https://github.com/chrisliu430)\n信箱: heranchris0430@gmail.com\n開源連結: https://github.com/chrisliu430/Howger-Blessing-Telegrambot"
		case "audio":
			format = "mp3"
		case "video":
			format = "mp4"
		default:
			if update.Message == nil {
				continue
			} else {
				botReponseMessage(bot, update.Message.Text, update.Message.Chat.ID, format)
			}
		}
		if msg.Text != "" {
			bot.Send(msg)
		}
	}
}

func callGoogleAnalytics() {
	log.Println("Analysis")
	analyticURL := "https://www.google-analytics.com/collect"
	requestForm := url.Values{
		"v":   {"1"},
		"tid": {os.Getenv("TID")},
		"t":   {"event"},
		"ec":  {"Howger"},
		"ea":  {"Blessing"},
	}
	resp, err := http.PostForm(analyticURL, requestForm)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
}

func cralwerToGetVideo(text string, format string) map[string]string {
	var result map[string]string
	requestForm := url.Values{
		"text":   {text},
		"format": {format},
	}
	resp, err := http.PostForm("http://howfun.macs1207.info/api/media", requestForm)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	return result
}
