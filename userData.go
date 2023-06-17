package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)
func checkCacheDir() error{
	_, err := os.Stat("./cache")
	if err == nil { return nil }
	if os.IsNotExist(err){
		err = os.MkdirAll("./.cache", 0700)
		if err != nil{
			return err
		}
		return nil
	}
	return err
}

func checkUserData(userId int) (bool,error){
	if err := checkCacheDir(); err != nil{
		log.Println("Error in func checkCacheDir(): ", err)
	}
	userPath := "./.cache/" + strconv.Itoa(userId)
	_, err := os.Stat(userPath)
	if err == nil{
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist){
		return false, nil
	}
	return false, err
}

func getUserData(userId int) (UserData,error) {
	userPath := "./.cache/" + strconv.Itoa(userId)
	if fileExist, err := checkUserData(userId); err != nil{
		log.Println("Error in func checkUserPath(): ", err)
		return UserData{}, err
	}else if fileExist{
		file, err := os.Open(userPath)
		if err != nil{
			return UserData{}, err
		}
		defer file.Close()
		body, err := ioutil.ReadAll(file)

		var userData UserData
		if err := json.Unmarshal(body, &userData); err != nil{
			return UserData{}, err
		}
		return userData, nil
	}else{
		return UserData{}, nil
	}
}

func createUserData(userId int, data UserData) error{
	userPath := "./.cache/" + strconv.Itoa(userId)
	if _, err := checkUserData(userId); err != nil{
		log.Println("Error in func checkUserPath(): ", err)
		return err
	}else{
		file, err := os.Create(userPath)
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
			return err
		}
		defer file.Close()
		dataJson, err := json.MarshalIndent(data, "", "\t")
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
		}
		_, err = file.Write(dataJson)
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
		}
		return nil
	}
}

func strToUserData(str string) (UserData,error){
	str2 := strings.Split(str, "\n")
	if len(str2) != 3{
		return UserData{}, errors.New("Wrong format of userData")
	}
	var flag int = 0
	var userData UserData
	for i, l := range str2{
		switch i {
		case 0:
			if !strings.HasPrefix(l, "login: "){
				return UserData{}, errors.New("Wrong format of userData")
			}
			userData.Login = l[strings.Index(l," ") + 1:]
		case 1:
			if !strings.HasPrefix(l, "password: "){
				return UserData{}, errors.New("Wrong format of userData")
			}
			userData.Password = l[strings.Index(l," ") + 1:]
		case 2:
			if !strings.HasPrefix(l, "server: "){
				return UserData{}, errors.New("Wrong format of userData")
			}
			userData.Server = l[strings.Index(l," ") + 1:]
		}
		flag ++
		}
	return userData, nil
}

func (userData UserData)  String() string{
	return "\nТекущие данные:\n" +
		"login: " + userData.Login + "\n" +
		"password: " + userData.Password + "\n" +
		"server: " + userData.Server + "\n"
}

func getPathTemplate (path string) string{
	ans := ""
	path = path[1:]
	path2 := strings.Split(path, "/")
	n := len(path2)
	if n > 0{
		switch path2[0]{
		case "nodes":
			ans += "/nodes"
			if n > 1{
				ans += "/node"
				if n > 2{
					switch path2[2]{
					case "lxc":
						ans += "/lxc"
						if n > 3 {
							ans += "/vmid"
							if n > 4 {
								switch path2[4]{
								case "status":
									ans += "/status"
									if n > 5 {
										switch path2[5]{
										case "current":
											ans += "/current"
										case "reboot":
											ans += "/reboot"
										case "shutdown":
											ans += "/shutdown"
										case "start":
											ans += "/start"
										default:
											return ""
										}
									}
								default:
									return ""
								}
							}
						}
					case "status":
						ans += "/status"
						if n > 3 {
							switch path2[3] {
							case "reboot":
								ans += "/reboot"
							case "shutdown":
								ans += "/shutdown"
							case "current":
								ans += "/current"
							default:
								return ""
							}
						}
					default:
						return ""
					}
				}
			}
		default:
			return ""
		}
	}
	return ans
}