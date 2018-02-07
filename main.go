package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	mailgun "github.com/mailgun/mailgun-go"
)

func main() {
	domain := os.Getenv("DOMAIN_MAILGUN_SANDBOX")
	apiKeyPrivate := os.Getenv("PRIVATE_KEY_MAILGUN")
	apiKeyPublic := os.Getenv("PUBLIC_KEY_MAILGUN")
	mg := mailgun.NewMailgun(domain, apiKeyPrivate, apiKeyPublic)
	message := mg.NewMessage(
		"jt@20scoops.net",
		"Fancy subject!",
		"",
		"pondthaitay@gmail.com")
	html, err := ioutil.ReadFile("templates/inlined/alert.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	message.SetHtml(string(html))
	_, _, err = mg.Send(message)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("send email successfully")
	}
}
