// Copyright 2015 Blackhawk Network, Inc.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chimpmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TemplateInfoRequest struct {
	ApiKey       string `json:"key"`
	TemplateName string `json:"name"`
}
type TemplateInfo struct {
	Slug             string
	Name             string
	Labels           []string
	Code             string
	Subject          string
	FromEmail        string
	FromName         string
	Text             string
	PublishName      string
	PublishCode      string
	PublishSubject   string
	PublishFromEmail string
	PublishFromName  string
	PublishText      string
	PublishedAt      string
	CreatedAt        string
	UpdatedAt        string
}

func GetTemplateInfo() (err error) {
	r := TemplateInfoRequest{ApiKey: "53yx5-nHBEYqKlyf8zfk8g", TemplateName: "transactional-notification"}
	buf, err := json.Marshal(r)
	b := bytes.NewBuffer(buf)
	resp, err := http.Post("https://mandrillapp.com/api/1.0/templates/info.json", "text/json", b)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var templateInfo TemplateInfo
	err = json.Unmarshal(body, &templateInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("templateInfo: %+v", templateInfo)
	return err
}
