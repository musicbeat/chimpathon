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
	name 			string
	content			string
}

type To struct {
	email 			string
	name 			string
}

type Message struct {
	subject 		string
	from_email 		string
	to 				[]To
	important 		bool
	global_merge_vars []GlobalMergeVars
}

type SendTemplateRequest struct {
	key 			string
	message 		Message
	send_at			string
}

func SendTemplate() (err error) {
	g := make([]GlobalMergeVars, 3)

	g[0] = GlobalMergeVars{
		name: "SENDER_MESSAGE",
		content: "Happy Birthday",
	}

	g[1] = GlobalMergeVars{
		name: "EGIFT_ID",
		content: "https://www.google.com",
	}

	g[2] = GlobalMergeVars{
		name: "PARTNER",
		content: "GCM",
	}

	t := make([]To, 3)

	t[0] = To{
		email: "rahul.dabas@bhnetwork.com",
		name: "Rahul Dabas",
	}

	t[1] = To{
		email: "rdabas@nexient.com",
		name: "Rahul Dabas",
	}

	t[2] = To{
		email: "dougbusley@gmail.com",
		name: "Doug Busley",
	}

	m := Message{
		subject: "Subject",
		from_email: "craig.thomas@bhnetwork.com",
		to:t,
		important: true,
		global_merge_vars: g,
	}

	r := SendTemplateRequest{
		key: "53yx5-nHBEYqKlyf8zfk8g", 
		message:m,
		send_at:"2015-03-10T12:00:00",
	}
	log.Print(r)

	buf, err := json.Marshal(r)
	

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
