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
	botResponse(bot, updates)
}

func botResponse(bot *TelegramBotAPI.BotAPI, updates TelegramBotAPI.UpdatesChannel) {
	for update := range updates {
		callGoogleAnalytics()
		msg := TelegramBotAPI.NewMessage(update.Message.Chat.ID, "")
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg.Text = "直接輸入文字即可\n若有想建議的服務\n可以寄信至heranchris0430@gmail.com或至github上提出issue\nHow哥並無唸英文，所以可以打相似的音來讓HOW哥念:)"
			bot.Send(msg)
		default:
			result = cralwerToGetVideo(update.Message.Text)
			log.Println(result)
			if result["error"] != "" {
				msg.Text = "Howger沒念這個字哦"
				bot.Send(msg)
			} else {
				videoURL := "http://howfun.macs1207.info/video?v=" + result["video_id"]
				msg.Text = "這個影片的網址\n" + videoURL
				videoMsg := TelegramBotAPI.NewVideoShare(update.Message.Chat.ID, videoURL)
				bot.Send(msg)
				bot.Send(videoMsg)
			}
			result = { "": ""}
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
