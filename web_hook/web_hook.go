package web_hook

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func Post(webHookUrl string, toUsername string, fromUsername string, message string, iconEmoji string) {
	payload := fmt.Sprintf(`{"channel": "%s", "username": "%s", "text": "%s", "icon_emoji": "%s"}`, toUsername, fromUsername, message, iconEmoji)
	vs := url.Values{}
	vs.Set("payload", payload)
	r, e := http.PostForm(webHookUrl, vs)
	if e != nil {
		log.Fatalln(e)
	}
	log.Println(r)
}
