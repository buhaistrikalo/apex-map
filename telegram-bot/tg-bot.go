package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	tgToken := os.Getenv("BOT_TOKEN")
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
		log.Printf("Got update: %v", update.Message)
		if update.Message != nil { 
			command := update.Message.Command()
			if command == "map" {
				sendMap(bot, update.Message.Chat.ID)
			}
			if command == "ping" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			sendMap(bot, chatID)
		}
	}
}

func sendMap(bot *tgbotapi.BotAPI, chatID int64) {
	var response, err = getMap()
	if err != nil {
		log.Println("Failed to get map:", err)
		return
	}
	var mapName string = response.Map
	var duration string = response.Duration
	var nextMap string = response.NextMap
	var nextMapTime string = response.NextMapTime

	var mapImage string = getMapImage(mapName)
	file, err := os.Open(mapImage)
	if err != nil {
		log.Println("Failed to open image:", err)
		return
	}
	reader := tgbotapi.FileReader{Name: "image.jpg", Reader: file}
	msg := tgbotapi.NewPhoto(chatID, reader)
	msg.Caption = "🖼 Карта: " + mapName + "\n⌛ Осталось: " + duration + "\n\nСледующая карта: " + nextMap + " - " + translateTime(nextMapTime)
	msg.ReplyMarkup = numericKeyboard
	bot.Send(msg)
}

func getMapImage(mapName string) string {
	var mapImage string
	if mapName == "Storm Point" {
		mapImage = "./images/storm_point.webp"
	} else if mapName == "World's Edge" {
		mapImage = "./images/worlds_edge.webp"
	} else if mapName == "Broken Moon" {
		mapImage = "./images/broken_moon.webp"
	} else {
		mapImage = "./images/olympus.webp"
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

func getMap() (*Response, error) {
	serverUrl := os.Getenv("SERVER_URL")
	url := serverUrl + "api/apex-map"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
		fmt.Println("Failed to get map: %v\n", err)
        return nil, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println("Failed to get map 1: %v\n", err)
        return nil, err
    }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to get map 2: %v\n", err)
		return nil, err
	}

	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Failed to get map 3: %v\n", err)
		return nil, err
	}

    return &response, nil
}