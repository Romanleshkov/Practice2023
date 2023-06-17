package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main(){
	if err := godotenv.Load("config.env"); err != nil{
		log.Fatal("file is not found: config.env")
	}
	// https://api.telegram.org/bot<token>/METHOD_NAME
	botToken := os.Getenv("BOT_TOKEN")
	botUrl := "https://api.telegram.org/bot" + botToken
	offset := 0
	for true{
		updates, err := getUpdates(botUrl, offset)
		if err != nil{
			log.Println("Error from func getUpdates: ", err.Error())
		}
		for _, update := range updates{
			//fmt.Println("update: ", update)
			if err := respondChat(botUrl, update); err != nil{
				log.Println("Error from func respondChat: ", err.Error())
			}
			if err := editMessage(botUrl, update); err != nil{
				log.Println("Error from func respondQuery: ", err.Error())
			}
			offset = update.UpdateId + 1
		}
	}

}

func getUpdates(botUrl string, offset int) ([]Update, error){
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var response Response
	if err := json.Unmarshal(body, &response); err != nil{
		return nil, err
	}
	return response.Result, nil
}

func respondChat(botUrl string, update Update) error {
	if update.Message.Chat.Id == 0{
		return nil
	}
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.Id
	switch update.Message.Text {
	case "/start":
		botMessage.Text = userDataMenu.text
		botMessage.ReplyMarkup.InlineKeyboard = userDataMenu.keyboard
	default:
		userData, err := strToUserData(update.Message.Text)
		if err != nil{
			botMessage.Text = "Неверный формат пользовательских данных\n" + userDataMenu.text
			botMessage.ReplyMarkup.InlineKeyboard = userDataMenu.keyboard
		}else {
			if err := createUserData(update.Message.From.Id, userData); err != nil{
				log.Println("Error in func createUserData(): ", err)
			}
			botMessage.Text = mainMenu.text + "\nДанные внесены.\n" + userData.String()
			botMessage.ReplyMarkup.InlineKeyboard = mainMenu.keyboard
		}
	}
	buf, err := json.Marshal(botMessage)
	if err != nil{
		return err
	}
	_, err = http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil{
		return err
	}
	return nil
}

func editMessage (botUrl string, update Update) error {
	if update.CallbackQuery.Message.Chat.Id == 0 {
		return nil
	}
	var answerQuery AnswerCallbackQuery
	answerQuery.CallbackQueryId = update.CallbackQuery.Id
	answerQuery.Text = ""
	var editMessage EditMessageText
	editMessage.ChatId = update.CallbackQuery.Message.Chat.Id
	editMessage.MessageId = update.CallbackQuery.Message.MessageId
	path := update.CallbackQuery.Data
	pathT := getPathTemplate(path)
	switch path {
	case "/mainMenu":
		editMessage.Text = "Главное меню"
		editMessage.ReplyMarkup.InlineKeyboard = mainMenu.keyboard
	case "/userDataMenu":
		editMessage.Text = userDataMenu.text
		if fileExist, err := checkUserData(update.CallbackQuery.From.Id); err != nil {
			log.Println("Error in func checkUserData(): ", err)
		} else if fileExist {
			userData, err := getUserData(update.CallbackQuery.From.Id)
			if err != nil {
				log.Println("Error in func getUserData(): ", err)
			}
			editMessage.Text = editMessage.Text + userData.String()
		}
		editMessage.ReplyMarkup.InlineKeyboard = userDataMenu.keyboard
	case "/actionMenu":
		editMessage.Text = actionMenu.text
		editMessage.ReplyMarkup.InlineKeyboard = actionMenu.keyboard
	}
	switch pathT {
	case "/nodes":
		answerJson, err := getUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func getUserRequest(): ", err)
		}
		var response GetNodes
		if err := json.Unmarshal(answerJson, &response); err != nil{
			log.Println("Error in func json.Unmarshal(): ", err)
		}
		menu := nodesMenu(response)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node":
		menu := nodeMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node/lxc":
		answerJson, err := getUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func getUserRequest(): ", err)
		}
		var response GetLxcs
		if err := json.Unmarshal(answerJson, &response); err != nil{
			log.Println("Error in func json.Unmarshal(): ", err)
		}
		menu := lxcsMenu(response, path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node/lxc/vmid":
		menu := lxcMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node/lxc/vmid/status":
		menu := lxcStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node/lxc/vmid/status/current":
		var response GetLxcStatus
		answerJson, err := getUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func getUserRequest(): ", err)
		}
		if err := json.Unmarshal(answerJson, &response); err != nil{
			log.Println("Error in func json.Unmarshal(): ", err)
		}
		path = path[:strings.LastIndex(path, "/")]
		menu := lxcStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = response.Data.Status
	case "/nodes/node/lxc/vmid/status/start":
		_, err := postUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func postUserRequest(): ", err)
		}
		path = path[:strings.LastIndex(path, "/")]
		menu := lxcStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = "Done"
	case "/nodes/node/lxc/vmid/status/shutdown":
		_, err := postUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func postUserRequest(): ", err)
		}
		path = path[:strings.LastIndex(path, "/")]
		menu := lxcStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = "Done"
	case "/nodes/node/lxc/vmid/status/reboot":
		_, err := postUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func postUserRequest(): ", err)
		}
		path = path[:strings.LastIndex(path, "/")]
		menu := lxcStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = "Done"
	case "/nodes/node/status":
		menu := nodeStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
	case "/nodes/node/status/current":
		var response GetNodeStatus
		path = path[:strings.LastIndex(path, "/")]
		answerJson, err := getUserRequests(update.CallbackQuery.From.Id, path)
		if err != nil{
			log.Println("Error in func getUserRequest(): ", err)
		}
		if err := json.Unmarshal(answerJson, &response); err != nil{
			log.Println("Error in func json.Unmarshal(): ", err)
		}
		menu := nodeStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		if response.Data.Uptime > 0{
			answerQuery.Text = "running"
		}else {
			answerQuery.Text = "stopped"
		}
	case "/nodes/node/status/shutdown":
		path = path[:strings.LastIndex(path, "/")]
		dataJS := []byte(`{"command":"shutdown"}`)
		_, err := postUserRequests(update.CallbackQuery.From.Id, path, dataJS)
		if err != nil{
			log.Println("Error in func postUserRequest(): ", err)
		}
		menu := nodeStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = "Done"
	case "/nodes/node/status/reboot":
		path = path[:strings.LastIndex(path, "/")]
		dataJS := []byte(`{"command":"reboot"}`)
		fmt.Println(dataJS)
		_, err := postUserRequests(update.CallbackQuery.From.Id, path, dataJS)
		if err != nil{
			log.Println("Error in func postUserRequest(): ", err)
		}
		menu := nodeStatusMenu(path)
		editMessage.Text = menu.text
		editMessage.ReplyMarkup.InlineKeyboard = menu.keyboard
		answerQuery.Text = "Done"
	}
	buf, err := json.Marshal(editMessage)
	if err != nil{
		return err
	}
	_, err = http.Post(botUrl + "/editMessageText", "application/json",
		bytes.NewBuffer(buf))
	if err != nil{
		return err
	}
	buf2, err := json.Marshal(answerQuery)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl + "/answerCallbackQuery", "application/json",
		bytes.NewBuffer(buf2))
	return nil
}
