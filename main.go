package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	// admins        []string
	telegramToken string

	login    string
	password string
	token    string

	timezone int

	commandKeyboard tgbotapi.ReplyKeyboardMarkup
)

func main() {
	flag.StringVar(&login, "login", "", "Mishiko login/email")
	flag.StringVar(&password, "password", "", "Mishiko password")
	flag.StringVar(&token, "token", "", "Mishiko token")
	flag.StringVar(&telegramToken, "telegram", "", "Telegram token")
	flag.IntVar(&timezone, "timezone", 3, "timezone")

	flag.Parse()

	if token == "" && login == "" && password == "" {
		log.Println("You should provide token or login/password to start")
	} else {
		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = false
		log.Printf("Authorized on account %s", bot.Self.UserName)

		var ucfg = tgbotapi.NewUpdate(0)
		ucfg.Timeout = 60

		updates, err := bot.GetUpdatesChan(ucfg)

		if err != nil {
			log.Fatalf("[INIT] [Failed to init Telegram updates chan: %v]", err)
		}

		for {
			select {
			case update := <-updates:
				// if intInStringSlice(int(update.Message.From.ID), admins) {

				// Text := update.Message.Text
				// Command := update.Message.Command()
				// Args := update.Message.CommandArguments()

				// msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

				sendStats(bot, int64(update.Message.From.ID))
				// }
			}
		}
	}
}

func sendStats(bot *tgbotapi.BotAPI, userID int64) {
	msg := tgbotapi.NewMessage(userID, "")

	result := ""
	pets := getPets(false)
	if len(pets) > 0 {
		for index := range pets {
			petData := pets[index]
			petActivity := getActivity(petData.ID, false)
			charging := ""
			if petActivity.Charging {
				charging = " (charging)"
				petActivity.BatteryCharge = int(math.Abs(float64(petActivity.BatteryCharge)))
			}
			// result = fmt.Sprintf("PetID: %d\nSteps: %d\nActivity: %d/%d\nDistance: %.3fm\nBattery: %d%%", petActivity.PetID, petActivity.CurrentEnergy, petActivity.CurrentActivity, petActivity.PetActivityAim, float64(petActivity.CurrentDistance)/1000, petActivity.BatteryCharge)
			result = fmt.Sprintf("Activity: %d/%d\nSteps: %d\nDistance: %.3fm\nBattery: %d%%%s", petActivity.CurrentActivity, petActivity.PetActivityAim, petActivity.CurrentEnergy, float64(petActivity.CurrentDistance)/1000, petActivity.BatteryCharge, charging)
		}
	}

	petsLocation := getPetsLocations(false)
	if len(petsLocation.Pets) > 0 {
		for index := range petsLocation.Pets {
			petData := petsLocation.Pets[index]

			// result += fmt.Sprintf("\nPetID: %d\nLat: %.6f\nLon: %.6f\nAlt: %.2f\nAccuracy: %.2f\nDate: %s\nSos: %d", petData.ID, petData.Lat, petData.Lon, petData.Alt, petData.Accuracy, time.Unix(petData.Date/1000, 0), petData.SosModeTime)
			timeData := time.Unix(petData.Date/1000, 0)
			result += fmt.Sprintf("\nLocation: %.6f, %.6f (±%.2fm)\nUpdated: %s", petData.Lat, petData.Lon, petData.Accuracy, timeData.UTC().Format("01.14.2006 15:04:05"))

			if petData.Lat != 0.0 && petData.Lon != 0.0 {
				location := tgbotapi.NewLocation(userID, petData.Lat, petData.Lon)
				bot.Send(location)
			}
		}
	}

	msg.Text = result

	// if userID == -1 {
	// 	for _, id := range admins {
	// 		userID, _ = strconv.ParseInt(id, 10, 64)
	// 		msg.ChatID = userID
	// 		bot.Send(msg)
	// 	}
	// } else {
	bot.Send(msg)

	// }
}

// func intInStringSlice(a int, list []string) bool {
// 	b := strconv.Itoa(a)
// 	for _, c := range list {
// 		if b == c {
// 			return true
// 		}
// 	}
// 	return false
// }
