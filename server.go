package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
)

type LineOfLog struct {
	RemoteAddr  string
	ContentType string
	Path        string
	Query       string
	Method      string
	Body        string
}

var TemplateOfLog = `
Remote address:   {{.RemoteAddr}}
Content-Type:     {{.ContentType}}
HTTP method:      {{.Method}}

path:
{{.Path}}

query string:
{{.Query}}

body:             
{{.Body}}

`

func main() {

	// Open HTML.
	http.Handle("/", http.FileServer(http.Dir(".")))

	// send mail to SMTP server.
	http.HandleFunc("/send", sendMailHandle)

	// Listen start 8080 port.
	http.ListenAndServe(":8080", http.DefaultServeMux)

}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		line := LineOfLog{
			r.RemoteAddr,
			r.Header.Get("Content-Type"),
			r.URL.Path,
			r.URL.RawQuery,
			r.Method, body,
		}
		tmpl, err := template.New("line").Parse(TemplateOfLog)
		if err != nil {
			panic(err)
		}

		bufline := new(bytes.Buffer)
		err = tmpl.Execute(bufline, line)
		if err != nil {
			panic(err)
		}

		log.Printf(bufline.String())
		handler.ServeHTTP(w, r)
	})
}

func sendMailHandle(w http.ResponseWriter, r *http.Request) {

	// Bodyデータを扱う場合には、事前にパースを行う
	r.ParseForm()

	// Get PostForm data.
	postForm := r.PostForm

	from := mail.Address{"", postForm.Get("from")}
	to := mail.Address{"", postForm.Get("to")}
	subject := postForm.Get("subject")
	password := postForm.Get("password")
	body := postForm.Get("body")

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
	message += "\r\n" + body

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
