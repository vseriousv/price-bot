package telegram

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/models"
	"github.com/vseriousv/price-bot/internal/price_alerts"
	"github.com/vseriousv/price-bot/internal/providers"
	"github.com/vseriousv/price-bot/internal/users"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func StartBot(db *pgxpool.Pool, token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// black list
		if getBlackList()[update.Message.Chat.ID] {
			defaultMessage(bot, &update, "Нет доступа")
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			commandPrice := strings.Split(update.Message.Text, " ")
			command := strings.Split(commandPrice[0], "__")

			switch command[0] {
			case "/start":
				getStartMessage(bot, &update, db)
			//case "/get_tickers":
			//	if len(command) > 1 {
			//		currency := command[1]
			//		getTickers(bot, &update, currency)
			//	} else {
			//		defaultMessage(bot, &update)
			//	}
			case "/get_price":
				if len(command) > 1 {
					ticker := strings.Join(strings.Split(command[1], "_"), "-")
					getPrice(bot, &update, ticker)
				} else {
					defaultMessage(bot, &update, "")
				}
			case "/create_alert":
				if len(command) > 1 && len(commandPrice) > 1 {
					ticker := strings.Join(strings.Split(command[1], "_"), "-")
					price := commandPrice[1]
					createAlert(bot, &update, db, ticker, price)
				} else {
					defaultMessage(bot, &update, "укажите цену через пробел")
				}
			case "/get_alerts":
				getAlerts(bot, &update, db)
			case "/delete_alert":
				alertId, err := strconv.Atoi(command[1])
				if err != nil {
					defaultMessage(bot, &update, "Ошибка удаления")
				} else {
					deleteAlert(bot, &update, db, int64(alertId))
					getAlerts(bot, &update, db)
				}
			default:
				defaultMessage(bot, &update, "")
			}
			//clearMessage(bot, &update)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
	return nil
}

func getStartMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, db *pgxpool.Pool) {
	text := `
👋 Привет! Я Price Monitoring Bot!
Чтобы запросить текущую цену используй следующую конструкцию команды:
▶️ /get_price__BTC_USDT

Если хочешь установить оповещение на тикер с определенной ценой,
тогда используй команду ниже и укажи цену через пробел:
▶️ /create_alert__BTC_USDT 16590.5

Когда цена будет достигнута, тебе предет оповещение в телеграм 👌

Посмотреть созданные оповещания:
▶️ /get_alerts
Вернет вам список в формате:
537 [BTC-USDT] 16590.50⬇ (16803.20)
Первая цифра - это порядковый номер
В скобках указана текущая цена на момент создания

Удалить оповещание (например: с порядковым номером 537):
▶️ /delete_alert__537
`
	// check exist user by chat_id
	check, err := users.GetByChatId(db, update.Message.Chat.ID)
	if err != nil {
		log.Println("[GetByChatId] :: ", err)
	}

	if check == nil {
		var user models.User
		user.ChatId = update.Message.Chat.ID
		user.UserName = update.Message.Chat.UserName
		user.FirstName = update.Message.Chat.FirstName
		user.LastName = update.Message.Chat.LastName
		user.Description = update.Message.Chat.Description
		//user.Photo = update.Message.Chat.Photo.BigFileID
		user.Title = update.Message.Chat.Title
		user.AllMembersAreAdmins = update.Message.Chat.AllMembersAreAdmins
		user.InviteLink = update.Message.Chat.InviteLink

		//create user
		log.Printf("[USER/CREATE] :: %d", user.ChatId)
		if err := user.Create(db); err != nil {
			text = "Не удалось авторизовать ползователя"
			log.Println("err", err)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func getTickers(bot *tgbotapi.BotAPI, update *tgbotapi.Update, currency string) {
	p, err := providers.GetProvider("kucoin")
	if err != nil {
		log.Println(err)
	}

	text := ""
	for _, ticker := range p.GetTickersList() {
		symbols := strings.Split(string(ticker), "-")
		cur := strings.ToUpper(currency)
		if symbols[0] == cur || symbols[1] == cur {
			text += fmt.Sprintf("/get_price__%s_%s\n", symbols[0], symbols[1])
		}

		//if i%10 == 0 && i != 0 {
		//	strings.ToUpper(text)
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		//
		//	bot.Send(msg)
		//	text = ""
		//}
	}
}

func getPrice(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ticker string) {
	p, err := providers.GetProvider("kucoin")
	if err != nil {
		log.Println(err)
	}

	text := fmt.Sprintf("[GET PRICE] :: %s", string(*p.GetPriceByTicker(ticker)))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func createAlert(bot *tgbotapi.BotAPI, update *tgbotapi.Update, db *pgxpool.Pool, ticker, price string) {
	// check exist user by chat_id
	u, err := users.GetByChatId(db, update.Message.Chat.ID)
	if err != nil {
		log.Println("[GetByChatId] :: ", err)
	}

	p, err := providers.GetProvider("kucoin")
	if err != nil {
		log.Println(err)
	}

	currentPrice, err := strconv.ParseFloat(string(*p.GetPriceByTicker(ticker)), 64)
	if err != nil {
		log.Println(err)
	}

	alertPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println(err)
	}

	var arrow string
	if currentPrice < alertPrice {
		arrow = "⬆"
	} else {
		arrow = "⬇"
	}
	text := fmt.Sprintf("[%s] %.8g%s (%.8g)", ticker, alertPrice, arrow, currentPrice)

	if u == nil {
		text = "Не удалось пройти авторизацию"
	} else {
		var pa models.PriceAlert
		pa.User = *u
		pa.Ticker = ticker
		pa.CreatePrice = currentPrice
		pa.AlertPrice = alertPrice
		err := pa.Create(db)
		if err != nil {
			log.Println(err)
			text = "Не удалось установить данную цену, попробуйте другую или позднее"
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func getAlerts(bot *tgbotapi.BotAPI, update *tgbotapi.Update, db *pgxpool.Pool) {
	var text = ""
	alerts, err := price_alerts.GetListByChatId(db, update.Message.Chat.ID)
	if err != nil || alerts == nil || len(*alerts) == 0 {
		text = "У вас нет созданных оповещений"
		log.Println("GetListByChatId", err)
	}

	for _, alert := range *alerts {
		var arrow string
		if alert.IsUp {
			arrow = "⬆"
		} else {
			arrow = "⬇"
		}
		text += fmt.Sprintf("%d [%s] %.8g%s (%.8g)\n", alert.Id, alert.Ticker, alert.AlertPrice, arrow, alert.CreatePrice)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func deleteAlert(bot *tgbotapi.BotAPI, update *tgbotapi.Update, db *pgxpool.Pool, id int64) {
	var text = ""

	alert, err := price_alerts.GetById(db, update.Message.Chat.ID, id)
	if err != nil {
		log.Println("GetById", err)
	}
	if alert == nil {
		text = "У вас нет созданных оповещений с этим ID"
		log.Println("GetById", err)
	} else {
		isRemoved, err := price_alerts.DeleteAlertById(db, id, alert.User.Id)
		if err != nil && !isRemoved {
			text = "Ошибка удаления, попробуйте позднее"
		}

		text = "Оповещение удалено"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func defaultMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, text string) {
	if text == "" {
		text = "I don't know the command"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func clearMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	_, err := bot.DeleteMessage(
		tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID - 1,
		})
	_, err = bot.DeleteMessage(
		tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID - 2,
		})
	if err != nil {
		log.Println(err)
	}
}
