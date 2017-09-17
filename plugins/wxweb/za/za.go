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
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

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

func za_giphy(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !strings.HasPrefix(msg.Content, ZA_CMD_PREFIX) {
		return
	}
	content := strings.TrimPrefix(msg.Content, ZA_CMD_PREFIX)

	search, err := http.Get("https://giphy.com/search/" + url.QueryEscape(content))
	if err != nil {
		logs.Error(err)
		return
	}
	defer search.Body.Close()
	searchBody, err := ioutil.ReadAll(search.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	// Crappy RE for extracting image source from search.Body.
	// One caveat, the WeChat gif upload limit is 1MiB.  So the best
	// effor is to use the 'downsampled' gif.
	crappyRE := regexp.MustCompile(`"fixed_width_downsampled": {"url": "([^"']+)",`)
	matches := crappyRE.FindAll(searchBody, -1)
	if len(matches) < 1 {
		logs.Error("cannot get giphy image links")
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gifUrl := crappyRE.FindStringSubmatch(string(matches[r.Intn(len(matches))%20][:]))[1]
	logs.Info("gifUrl: ", gifUrl)
	resp, err := http.Get(gifUrl)
	if err != nil {
		logs.Error(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	// XXX. Report sending error, normally caused by sending larger than 1MiB data.
	if msg.FromUserName == session.Bot.UserName {
		session.SendEmotionFromBytes(b, session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendEmotionFromBytes(b, session.Bot.UserName, msg.FromUserName)
	}
}
