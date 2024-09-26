package email

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, matchName string, homePlayerScore int, awayPlayerScore int, status string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Match Status Update - "+matchName)
	m.SetBody("text/html", "<p>Match "+matchName+" is now <strong>"+status+" </strong></p><br/><p>Match Result :"+strconv.Itoa(homePlayerScore)+" - "+strconv.Itoa(awayPlayerScore)+"</p>")

	d := gomail.NewDialer("10.133.66.119", 1025, "snooker-data@local.dev", "")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
	}
}
