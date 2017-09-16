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

package jav

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)

var (
	JAV_CMD_PREFIX = "jav"
)

// register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(jav), "jav", JAV_CMD_PREFIX)
	if err := session.HandlerRegister.EnableByName("jav"); err != nil {
		logs.Error(err)
	}
}

func jav(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !strings.HasPrefix(msg.Content, JAV_CMD_PREFIX) {
		return
	}
	uri := "http://www.javlibrary.com/cn/vl_newrelease.php"
	s, err := spider.CreateSpiderFromUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	srcs, _ := s.GetAttr("div.videothumblist>div.videos>div.video>a>img", "src")
	if len(srcs) < 1 {
		logs.Error("cannot get most wanted JAV ids")
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Jav library always returns 20 results in one page.
	imgUrl := "http:" + strings.ToLower(strings.Replace(srcs[r.Intn(len(srcs))%20], "ps.jpg", "pl.jpg", 1))
	logs.Info("jav imgUrl: ", imgUrl)
	resp, err := http.Get(imgUrl)
	if err != nil {
		logs.Error(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}

	session.SendImgFromBytes(b, imgUrl, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
}
