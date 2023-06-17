package main

type UserData struct {
	Login string		`json:"login"`
	Password string		`json:"password"`
	Server string		`json:"server"`
}

type menu struct {
	keyboard [][]InlineKeyboardButton
	text string
}
