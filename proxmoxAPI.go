package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func getAccessTicket(username string, password string, serverIP string) (AccessTicketResponse, error) {
	var accessTicket AccessTicket
	accessTicket.Username = username
	accessTicket.Password = password
	proxmoxURL := "https://" + serverIP + "/api2/json"
	buf, err := json.Marshal(accessTicket)
	if err != nil{
		log.Println("Error in func json.Marshal(): ", err)
		return AccessTicketResponse{}, err
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(proxmoxURL + "/access/ticket", "application/json", bytes.NewBuffer(buf))
	if err != nil{
		log.Println("Error in func http.Post(): ", err)
		return AccessTicketResponse{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("Error in func ioutil.ReadAll(): ", err)
		return AccessTicketResponse{}, err
	}
	var response AccessTicketResponse
	if err := json.Unmarshal(body, &response); err != nil{
		log.Println("Error in func json.Unmarshal(): ", err)
		return AccessTicketResponse{}, err
	}
	if response.Data == (AccessTicketResponse{}).Data{
		return AccessTicketResponse{}, errors.New("wrong userData")
	}
	return response, nil
}

func getUserRequests (userId int, path string) ([]byte, error){
	userData, err := getUserData(userId)
	if err != nil{
		log.Println("Error in func getUserData()")
	}
	accessTicket, err := getAccessTicket(userData.Login, userData.Password, userData.Server)
	if err != nil{
		return nil, err
	}
	proxmoxURL := "https://" + userData.Server + "/api2/json"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", proxmoxURL + path, nil)
	if err != nil{
		log.Println("Error in func http.NewRequest(): ", err)
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: accessTicket.Data.Ticket})
	resp, err := client.Do(req)
	if err != nil{
		log.Println("Error in func client.Do(): ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("Error in func ioutil.ReadAll(): ", err)
		return nil, err
	}
	return body, nil
}

func postUserRequests (userId int, path string, jsData ...[]byte) ([]byte, error){
	userData, err := getUserData(userId)
	if err != nil{
		log.Println("Error in func getUserData()")
	}
	accessTicket, err := getAccessTicket(userData.Login, userData.Password, userData.Server)
	if err != nil{
		return nil, err
	}
	proxmoxURL := "https://" + userData.Server + "/api2/json"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	var dataJS []byte = nil
	if len(jsData) > 0{
		dataJS = jsData[0]
	}
	req, err := http.NewRequest("POST", proxmoxURL + path, bytes.NewReader(dataJS))
	if err != nil{
		log.Println("Error in func http.NewRequest(): ", err)
		return nil, err
	}
	if len(jsData) > 0{
		req.Header.Set("Content-Type", "application/json")
	}
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: accessTicket.Data.Ticket})
	req.Header.Add("CSRFPreventionToken", accessTicket.Data.CSRFPreventionToken)
	resp, err := client.Do(req)
	if err != nil{
		log.Println("Error in func client.Do(): ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("Error in func ioutil.ReadAll(): ", err)
		return nil, err
	}
	return body, nil
}
