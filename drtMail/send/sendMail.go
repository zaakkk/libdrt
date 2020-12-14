package send

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/zaakkk/libdrt/drtMail/coreMail"
)

var mailQueue chan coreMail.MailStruct = make(chan coreMail.MailStruct, 300)

// SendMailToSMTP はサーバが受け取ったメール群を順番にSMTPサーバへ送信する
func SendMailToSMTP(m coreMail.MailStruct) {
	headers := make(map[string]string)
	headers["From"] = m.From
	headers["To"] = m.To
	headers["Subject"] = m.Sub

	//from := mail.Address{"", m.From}
	//to := mail.Address{"", m.To}

	message := ""
	//add headers
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// add message
	message += "\r\n" + string(m.Msg)

	// Connect to the SMTP Server
	//servername := "smtp.mail.yahoo.co.jp:465"
	servername := "smtp.gmail.com:465"
	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", m.From, m.Password, host)

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
	//if err = c.Mail(from.Address); err != nil {
	if err = c.Mail(m.From); err != nil {
		log.Panic(err)
	}
	//if err = c.Rcpt(to.Address); err != nil {
	if err = c.Rcpt(m.To); err != nil {
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

	//w.Write([]byte(message))

}

//GMailSend is send mail parameter to Web server.
func GMailSend(m coreMail.MailStruct) error {

	//適宜URLを変更する(ec2/local)
	//urlTarget := "http://***.***.***.***:80/send"
	urlTarget := "http://localhost:8080/send"

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

// SendMailHandler はPOST処理されたメール群をキューに格納する
/*
func SendMailHandle(w http.ResponseWriter, r *http.Request) {

	// Bodyデータを扱う場合には、事前にパースを行う
	r.ParseForm()

	// Get PostForm data.
	postForm := r.PostForm

	//from := mail.Address{"", postForm.Get("from")}
	from := postForm.Get("from")
	//to := mail.Address{"", postForm.Get("to")}
	to := postForm.Get("to")
	subject := postForm.Get("subject")
	password := postForm.Get("password")
	body := postForm.Get("body")

	mailStruct := coreMail.MailStruct{
		From:     from,         // from
		Username: from,         // username
		Password: password,     // password
		To:       to,           // to
		Sub:      subject,      // subject
		Msg:      []byte(body), // message
	}
	fmt.Println("Subject: " + mailStruct.Sub)
	go func() {
		mailQueue <- mailStruct
	}()
	fmt.Println("len mailQueue: %d", len(mailQueue))
	//http.Redirect(w, r, "http://localhost:8080/sendMailToSMTP", http.StatusMovedPermanently)
	//<-mailQueue

	// TODO: 200返す
}
*/

// webサーバからAPIサーバへ断片を送信
func SendMailHandle(w http.ResponseWriter, r *http.Request) {

	// Bodyデータを扱う場合には、事前にパースを行う
	r.ParseForm()

	// Get PostForm data.
	postForm := r.PostForm

	//from := mail.Address{"", postForm.Get("from")}
	from := postForm.Get("from")
	//to := mail.Address{"", postForm.Get("to")}
	to := postForm.Get("to")
	subject := postForm.Get("subject")
	password := postForm.Get("password")
	body := postForm.Get("body")

	mailStruct := coreMail.MailStruct{
		From:     from,         // from
		Username: from,         // username
		Password: password,     // password
		To:       to,           // to
		Sub:      subject,      // subject
		Msg:      []byte(body), // message
	}

	mToJson, _ := json.Marshal(mailStruct)
	//fmt.Println(string(mToJson))

	conn, err := net.Dial("tcp", "localhost:50030")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, string(mToJson))
}
