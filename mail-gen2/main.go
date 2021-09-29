package main

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/appengine/v2/mail"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("index.html"))
		t.Execute(w, nil)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		// ローカルでは Metadata Server にアクセスできずクラッシュする
		// appengine.IsAppEngine() でローカル環境か判定できる
		log.Debugf(r.Context(), "subject=%s, body=%s", r.PostFormValue("subject"), r.PostFormValue("body"))

		msg := &mail.Message{
			Sender:  "anything@pokutuna-playground.appspotmail.com",
			To:      []string{"pokutuna <mail@pokutuna.com>"},
			Subject: r.PostFormValue("subject"),
			Body:    r.PostFormValue("body"),
		}
		if err := mail.Send(r.Context(), msg); err != nil {
			log.Errorf(r.Context(), "%v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	appengine.Main()
}
