package send

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"

	"../coreMail"
)

//GMailSend is send mail parameter to Web server.
func GMailSend(m coreMail.MailStruct) error {

	//適宜URLを変更する(ec2/local)
	urlTarget := "http://***.***.***.***:80/send"
	//urlTarget := "http://localhost:8080/send"

<<<<<<< HEAD
=======
	//urlTarget := "http://***.***.***.***:80/send"
	urlTarget := "http://localhost:8080/send"
>>>>>>> 113df24cf0935105cafe609375a8f67144c60d34
	args := url.Values{}
	args.Add("from", m.From)
	args.Add("to", m.To)
	args.Add("password", m.Password)
	args.Add("subject", m.Sub)
	args.Add("body", base64.StdEncoding.EncodeToString(m.Msg))

	// Debug
	//fmt.Println("args:\n" + m.From + "\n" + m.To + "\n" + m.Password + "\n" + m.Sub + "\n" + string(m.Msg) + "\n")

	res, err := http.PostForm(urlTarget, args)
	if err != nil {
		fmt.Println("Request error:", err)
		return err
	}
	defer res.Body.Close()

	// Check response
	//if err := checkResponse(res); err != nil {
	//	return err
	//}

	return nil
}

func checkResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		return err
	}

	strJSON := string(body)
	fmt.Println("checkResponse\n" + strJSON)
	return nil
}

// SendMailHandler is send mail from web server to SMTP server
func SendMailHandle(w http.ResponseWriter, r *http.Request) {

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
	//servername := "smtp.mail.yahoo.co.jp:465"
	servername := "smtp.gmail.com:465"

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

	w.Write([]byte(message))

}
