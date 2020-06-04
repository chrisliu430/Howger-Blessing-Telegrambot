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
	var bot *TelegramBotAPI.BotAPI
	var updates TelegramBotAPI.UpdatesChannel
	// --- Bot API Setting
	port := os.Getenv("PORT")
	url := os.Getenv("URL")
	botToken = os.Getenv("Token")
	addr := fmt.Sprintf(":%s", port)
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
	botProcessMessage(bot, updates)
}

func botReponseMessage(bot *TelegramBotAPI.BotAPI, text string, chatID int64) {
	msg := TelegramBotAPI.NewMessage(chatID, "")
	result := cralwerToGetVideo(text)
	log.Println(result)
	if result["error"] != "" {
		msg.Text = "Howger沒念這個字哦"
		bot.Send(msg)
	} else {
		videoURL := "http://howfun.macs1207.info/video?v=" + result["video_id"]
		msg.Text = "這個影片的網址\n" + videoURL
		videoMsg := TelegramBotAPI.NewVideoShare(chatID, videoURL)
		bot.Send(videoMsg)
	}
}

func botProcessMessage(bot *TelegramBotAPI.BotAPI, updates TelegramBotAPI.UpdatesChannel) {
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
			msg.Text = "--- Help List ---\n/start 可以查看如何使用\n/feedback 可以填寫建議(未完成)\n/voice [文字] 可以取得需要的音訊(未完成)\n/info 提供作者群、作者聯繫信箱及開源連結"
		case "info":
			msg.Text = "作者: C.H(https://github.com/chrisliu430)\n信箱: heranchris0430@gmail.com\n開源連結: https://github.com/chrisliu430/Howger-Blessing-Telegrambot"
		default:
			if update.Message == nil {
				continue
			} else {
				botReponseMessage(bot, update.Message.Text, update.Message.Chat.ID)
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

func cralwerToGetVideo(text string) map[string]string {
	var result map[string]string
	requestForm := url.Values{
		"text": {text},
	}
	resp, err := http.PostForm("http://howfun.macs1207.info/api/video", requestForm)
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
