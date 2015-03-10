// Copyright 2015 Blackhawk Network, Inc.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type Order struct {
	Title              string
	OrderNumber        string
	ProductId          string
	ProductDescription string
	FaceValue          string
	ToEmail            string
	EGiftURL           string
}

func (o *Order) save() error {
	// do something to email the order
	// save the order
	filename := o.Title + ".txt"
	j, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, j, 0600)
}

func loadOrder(title string) (*Order, error) {
	filename := title + ".txt"
	var order Order
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return &order, err
	}
	err = json.Unmarshal(b, &order)
	if err != nil {
		log.Fatal(err)
	}
	return &order, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	o, err := loadOrder(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+template.HTMLEscapeString(title), http.StatusFound)
		return
	}
	renderTemplate(w, "view", o)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	o, err := loadOrder(title)
	if err != nil {
		uuid := template.HTMLEscapeString(uuid.NewRandom().String())
		eGiftURL := "https://egift.io/" + uuid
		o = &Order{Title: uuid, OrderNumber: uuid, EGiftURL: eGiftURL}
	}
	renderTemplate(w, "edit", o)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	orderNumber := r.FormValue("orderNumber")
	productId := r.FormValue("productId")
	productDescription := r.FormValue("productDescription")
	faceValue := r.FormValue("faceValue")
	toEmail := r.FormValue("toEmail")
	eGiftURL := r.FormValue("eGiftURL")
	o := &Order{
        Title: title, 
        OrderNumber: orderNumber,
        ProductId: productId,
        ProductDescription: productDescription,
        FaceValue: faceValue,
        ToEmail: toEmail,
        EGiftURL: eGiftURL}
	err := o.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+template.HTMLEscapeString(title), http.StatusFound)
}

var funcMap = template.FuncMap{"htmlEscape": template.HTMLEscapeString}
var templates = template.Must(template.New("root").Funcs(funcMap).ParseFiles("edit.html","view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, o *Order) {
	err := templates.ExecuteTemplate(w, tmpl+".html", o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9-]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
