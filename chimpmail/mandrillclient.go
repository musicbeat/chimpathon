// Copyright 2015 Blackhawk Network, Inc.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chimpmail

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GlobalMergeVars struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type To struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Message struct {
	Subject         string            `json:"subject"`
	FromEmail       string            `json:"from_email"`
	To              []To              `json:"to"`
	Important       bool              `json:"important"`
	GlobalMergeVars []GlobalMergeVars `json:"global_merge_vars"`
}

type SendTemplateRequest struct {
	Key             string            `json:"key"`
	TemplateName    string            `json:"template_name"`
	Message         Message           `json:"message"`
	TemplateContent []GlobalMergeVars `json:"template_content"`
}

func SendTemplate(senderMessage string, eGiftId string, toEmail string, partner string) (err error) {
	g := make([]GlobalMergeVars, 3)

	g[0] = GlobalMergeVars{
		Name:    "SENDER_MESSAGE",
		Content: senderMessage,
	}

	g[1] = GlobalMergeVars{
		Name:    "EGIFT_ID",
		Content: eGiftId,
	}

	g[2] = GlobalMergeVars{
		Name:    "PARTNER",
		Content: partner,
	}

	t := make([]To, 3)

	t[0] = To{
		Email: toEmail,
		Name:  toEmail,
	}

	t[1] = To{
		Email: "rdabas@nexient.com",
		Name:  "Rahul Dabas",
	}

	t[2] = To{
		Email: "dougbusley@gmail.com",
		Name:  "Doug Busley",
	}

	m := Message{
		Subject:         senderMessage,
		FromEmail:       "craig.thomas@bhnetwork.com",
		To:              t,
		Important:       true,
		GlobalMergeVars: g,
	}

	r := SendTemplateRequest{
		Key:          "53yx5-nHBEYqKlyf8zfk8g",
		Message:      m,
		TemplateName: "transactional-notification",
	}
	log.Print(r)

	buf, err := json.Marshal(r)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(buf))

	b := bytes.NewBuffer(buf)
	resp, err := http.Post("https://mandrillapp.com/api/1.0/messages/send-template.json", "text/json", b)

	log.Print(resp.StatusCode)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Print(body)

	s := string(body)
	log.Print(s)

	return err
}
