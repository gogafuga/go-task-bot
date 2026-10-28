package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	urlpkg "net/url"
	"strconv"
)

func GetUpdate(token string, offset int) ([]Update, error) {
	baseURL := "https://api.telegram.org/bot" + token + "/getUpdates"
	params := urlpkg.Values{}
	params.Set("offset", strconv.Itoa(offset))
	params.Set("timeout", "30")
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response TelegramResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}

func SendMessage(token string, chatID int64, message string) error {
	baseURL := "https://api.telegram.org/bot" + token + "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("text", message)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
