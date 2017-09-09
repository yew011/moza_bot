/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package za

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)

var (
	ZA_CMD_PREFIX = "\u624e "
)

// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(za_giphy), "za", ZA_CMD_PREFIX)
	if err := session.HandlerRegister.EnableByName("za"); err != nil {
		logs.Error(err)
	}
}

func getGifUrl(url string) (string, error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return "", e
	}

	if res == nil {
		return "", errors.New("Response is nil")
	}

	defer res.Body.Close()
	if res.Request == nil {
		return "", errors.New("Response.Request is nil")
	}

	frame := map[string]interface{}{}
	json.NewDecoder(res.Body).Decode(&frame)
	jq := jsonq.NewQuery(frame)

	gifUrl, e := jq.String("data", "images", "original", "url")
	if e != nil {
		return "", e
	}

	return gifUrl, nil
}

func za_giphy(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !strings.HasPrefix(msg.Content, ZA_CMD_PREFIX) {
		return
	}
	content := strings.TrimPrefix(msg.Content, ZA_CMD_PREFIX)
	uri := "http://api.giphy.com/v1/gifs/translate?api_key=dc6zaTOxFJmzC&rating=r&s=" + url.QueryEscape(content)
	logs.Info("giphy translate url: %v", uri)
	gif, err := getGifUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	resp, err := http.Get(gif)
	if err != nil {
		logs.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if msg.FromUserName == session.Bot.UserName {
		session.SendEmotionFromBytes(body, session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendEmotionFromBytes(body, session.Bot.UserName, msg.FromUserName)
	}
}
