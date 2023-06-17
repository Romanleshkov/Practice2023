package main

import "strings"

var userDataMenu menu = menu{
	text: "Укажите логин, пароль и сервер\n" +
		"В формате:\n" +
		"login: <login>\n" +
		"password: <password>\n" +
		"server: <serverIP:port>\n",
	keyboard: [][]InlineKeyboardButton{
		{
			InlineKeyboardButton{Text: "Главное меню", CallbackData: "/mainMenu"},
		},
	},
}

var mainMenu menu = menu{
	text: "Главное меню",
	keyboard: [][]InlineKeyboardButton{
		{
			InlineKeyboardButton{Text: "Изменить данные",
				CallbackData: "/userDataMenu"},
			InlineKeyboardButton{Text: "Действия",
				CallbackData: "/actionMenu"},
		},
	},
}

var actionMenu menu = menu{
	text: "Действия",
	keyboard: [][]InlineKeyboardButton{
		{
			InlineKeyboardButton{Text: "Список узлов",
				CallbackData: "/nodes"},
		},
		{
			InlineKeyboardButton{Text: "Назад",
				CallbackData: "/mainMenu"},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	},
}

func nodesMenu(nodes GetNodes) menu{
	menu := menu{text: "Список узлов"}
	for _, node := range nodes.Nodes{
		menu.keyboard = append(
			menu.keyboard,
				[]InlineKeyboardButton{
					InlineKeyboardButton{Text: node.Node,
						CallbackData: "/nodes/" + node.Node},
				},
		)
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Действия",
				CallbackData: "/actionMenu"},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func nodeMenu(path string) menu{
	menu := menu{text: "Опции узла " + path}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Список lxc-контейнеров",
				CallbackData: path + "/lxc"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Статус узла",
				CallbackData: path + "/status"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Назад", CallbackData:
				path[:strings.LastIndex(path, "/")]},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func nodeStatusMenu(path string) menu{
	menu := menu{text: "Действия над статусом " +
		path[:strings.LastIndex(path, "/")]}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Текущий статус",
				CallbackData: path + "/current"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Перезагрузить",
				CallbackData: path + "/reboot"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Выключить",
				CallbackData: path + "/shutdown"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Назад", CallbackData:
			path[:strings.LastIndex(path, "/")]},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func lxcsMenu(lxcs GetLxcs, path string) menu{
	menu := menu{text: "Список lxc-контейнеров " +
		path[:strings.LastIndex(path, "/")]}
	for _, lxc := range lxcs.Lxcs{
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				InlineKeyboardButton{Text: lxc.VmId,
					CallbackData: path + "/" + lxc.VmId},
			},
		)
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Назад", CallbackData:
				path[:strings.LastIndex(path, "/")]},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func lxcMenu(path string) menu{
	menu := menu{text: "Опции lxc-контейнера " + path}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Статус",
				CallbackData: path + "/status"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Назад", CallbackData:
				path[:strings.LastIndex(path, "/")]},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func lxcStatusMenu(path string) menu{
	menu := menu{text: "Действия над статусом " +
		path[:strings.LastIndex(path, "/")]}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Текущий статус",
				CallbackData: path + "/current"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Перезагрузить",
				CallbackData: path + "/reboot"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Включить",
				CallbackData: path + "/start"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Выключить",
				CallbackData: path + "/shutdown"},
		},
		[]InlineKeyboardButton{
			InlineKeyboardButton{Text: "Назад",
				CallbackData:
				path[:strings.LastIndex(path, "/")]},
			InlineKeyboardButton{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}