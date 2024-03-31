package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Response struct {
	Map    string
	Duration string
	NextMap string
	NextMapTime string
	DateNow string
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Обновить", "/map"),
    ),
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tgToken := os.Getenv("TG_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { 
			if update.Message.Text == "/map" {
				sendMap(bot, update.Message.Chat.ID, update.Message.Date)
			}
		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			date := update.CallbackQuery.Message.Date
			resp, err := bot.Request(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
				continue
			}
			if !resp.Ok {
				log.Printf("Failed to delete message: response not OK")
				continue
			}
			sendMap(bot, chatID, date)
		}
	}
}

func sendMap(bot *tgbotapi.BotAPI, chatID int64, date int) {
	var response, err = getMap(int64(date))
	if err != nil {
		log.Println("Failed to get map:", err)
		return
	}
	var mapName string = response.Map
	var duration string = response.Duration
	var nextMap string = response.NextMap
	var nextMapTime string = response.NextMapTime

	var mapImage string = getMapImage(mapName)
	file, _ := os.Open(mapImage)

	reader := tgbotapi.FileReader{Name: "image.jpg", Reader: file}
	msg := tgbotapi.NewPhoto(chatID, reader)
	msg.Caption = "🖼 Карта: " + mapName + "\n⌛ Осталось: " + duration + "\n\nСледующая карта: " + nextMap + " - " + translateTime(nextMapTime)
	msg.ReplyMarkup = numericKeyboard
	bot.Send(msg)
}

func getMapImage(mapName string) string {
	var mapImage string
	if mapName == "Storm Point" {
		mapImage = "../image/storm_point.webp"
	} else if mapName == "World's Edge" {
		mapImage = "../image/worlds_edge.webp"
	} else if mapName == "Broken Moon" {
		mapImage = "../image/broken_moon.webp"
	} else {
		mapImage = "../image/olympus.webp"
	}
	return mapImage
}

func translateTime(inputTime string) string {
	t, err := time.Parse("15:04:05", inputTime)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	t = t.Add(time.Hour * 3)
	return t.Format("15:04:05")
}

func getMap(time int64) (*Response, error) {
	serverUrl := os.Getenv("SERVER_URL")
	url := serverUrl + "api/apex-map?time=" + strconv.FormatInt(time, 10)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

    return &response, nil
}