package main

import "net/smtp"

const (
	_from     = "sooot193@gmail.com"
	_smtpHost = "smtp.gmail.com"
	_smtpPort = "587"
)

var _mailAuth = smtp.PlainAuth("", _from, "oywtcjtnvsrhpbmp", _smtpHost)

func SendMail(toMail, msg string) error {
	return smtp.SendMail(_smtpHost+":"+_smtpPort, _mailAuth, _from, []string{toMail}, []byte(msg))
}
