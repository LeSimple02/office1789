package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type plainAuthInsecure struct {
	identity string
	username string
	password string
	host     string
}

func (a *plainAuthInsecure) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	resp := []byte("\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuthInsecure) Next(_ []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("unexpected server challenge")
	}
	return nil, nil
}

// sendSMTP opens a connection, negotiates STARTTLS if offered, authenticates if auth is provided,
// and sends the message. InsecureSkipVerify is enabled to allow self-signed certs in dev.
func sendSMTP(host, port, from string, to []string, msg []byte, auth smtp.Auth) error {
	addr := host + ":" + port
	hostForAuth := host
	if strings.Contains(hostForAuth, ":") {
		hostForAuth = strings.Split(hostForAuth, ":")[0]
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		cfg := &tls.Config{
			ServerName:         hostForAuth,
			InsecureSkipVerify: true, // dev convenience; mailserver uses self-signed
		}
		if err := c.StartTLS(cfg); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("auth: %w", err)
			}
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("rcpt to %s: %w", rcpt, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("write msg: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}
	return c.Quit()
}
