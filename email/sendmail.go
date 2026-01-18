package email

import (
	"log"
	"strconv"

	"github.com/astoyanov87/subscription-service/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, matchName string, homePlayerScore int, awayPlayerScore int, status string, cfg *config.Config) {
	m := gomail.NewMessage()
	m.SetHeader("From", "snooker-livescore@maildev.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Match Status Update - "+matchName)
	m.SetBody("text/html", "<p>Match "+matchName+" is now <strong>"+status+" </strong></p><br/><p>Match Result :"+strconv.Itoa(homePlayerScore)+" - "+strconv.Itoa(awayPlayerScore)+"</p>")

	d := gomail.NewDialer(cfg.SmtpConfig.Host, cfg.SmtpConfig.Port, cfg.SmtpConfig.Username, cfg.SmtpConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
	}
}
