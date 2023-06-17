package main

type Update struct {
	UpdateId int				`json:"update_id"`
	Message       Message       `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

type Message struct{
	MessageId int								`json:"message_id"`
	From User									`json:"from"`
	Chat struct {
		Id int									`json:"id"`
	}											`json:"chat"`
	Text string									`json:"text"`
	ReplyMarkup InlineKeyboardMarkup 			`json:"reply_markup"`
}

type CallbackQuery struct {
	Id string			`json:"id"`
	From User			`json:"from"`
	Data string			`json:"data"`
	Message Message		`json:"message"`
}

type User struct {
	Id int				`json:"id"`
	Username string		`json:"username"`
}

type Response struct {
	Result []Update 	`json:"result"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton 	`json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text string			`json:"text"`
	CallbackData string	`json:"callback_data"`
}

type BotMessage struct {
	ChatId	int									`json:"chat_id"`
	Text string									`json:"text"`
	ReplyMarkup InlineKeyboardMarkup 			`json:"reply_markup"`
}

type EditMessageText struct {
	ChatId int							`json:"chat_id"`
	MessageId int						`json:"message_id"`
	Text string							`json:"text"`
	ReplyMarkup InlineKeyboardMarkup	`json:"reply_markup"`
}

type AnswerCallbackQuery struct {
	CallbackQueryId string	`json:"callback_query_id"`
	Text string				`json:"text"`
}
