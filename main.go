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
	startText := "直接輸入文字即可\n若有想建議的服務\n可以寄信至heranchris0430@gmail.com或至github上提出issue\nHow哥並無唸英文，所以可以打相似的音來讓HOW哥念:)"
	for update := range updates {
		googleAnalytics()
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := TelegramBotAPI.NewMessage(update.Message.Chat.ID, startText)
			bot.Send(msg)
		default:
			requestForm := url.Values{
				"text": {update.Message.Text},
			}
			resp, err := http.PostForm("http://howfun.macs1207.info/api/video", requestForm)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(result)
			if result == nil {
				msg := TelegramBotAPI.NewVideoShare(update.Message.Chat.ID,
					"http://howfun.macs1207.info/video?v=1a9e30d6-6ddd-49d8-8d15-88d4a4f9d347")
				bot.Send(msg)
			} else {
				videoURL := "http://howfun.macs1207.info/video?v=" + result["video_id"]
				msg := TelegramBotAPI.NewVideoShare(update.Message.Chat.ID, videoURL)
				bot.Send(msg)
			}
		}
	}
	return
}

func googleAnalytics() {
	log.Println("Analysis")
	analyticURL := "https://www.google-analytics.com/collect"
	requestForm := url.Values{
		"v":   {"1"},
		"tid": {"UA-168546559-1"},
		"t":   {"event"},
		"ec":  {"Howger"},
		"ea":  {"Blessing"},
	}
	resp, err := http.PostForm(analyticURL, requestForm)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.StatusCode)
	defer resp.Body.Close()
	return
}
