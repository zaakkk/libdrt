package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
)

func main() {

	// Open HTML.
	http.Handle("/", http.FileServer(http.Dir(".")))

	// send mail to SMTP server.
	http.HandleFunc("/send", sendMailHandle)

	// Listen start 8080 port.
	http.ListenAndServe(":8080", nil)
}

func sendMailHandle(w http.ResponseWriter, r *http.Request) {

	// Bodyデータを扱う場合には、事前にパースを行う
	r.ParseForm()

	// Get PostForm data.
	form := r.PostForm
	//fmt.Fprintf(w, "フォーム：%v\n", form)

	// Get Form data
	//params := r.Form
	//fmt.Fprintf(w, "フォーム2：%v\n", params)

	from := mail.Address{"", form.Get("from")}
	to := mail.Address{"", form.Get("to")}
	subject := form.Get("subject")
	password := form.Get("password")
	body := form.Get("body")

	// Debug
	//fmt.Println("server:\n" + form.Get("from") + "\n" + form.Get("to") + "\n" + password + "\n" + subject + "\n" + body + "\n")

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// Change to Base64
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the SMTP Server
	servername := "smtp.mail.yahoo.co.jp:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", from.Address, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	data, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = data.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = data.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}
