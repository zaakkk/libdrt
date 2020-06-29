package recieve

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"../coreMail"

	"../../go-pop3"
)

//const yahooPOPServer = "pop.mail.yahoo.co.jp:995"
//const timeout = time.Second * 5

//sub(件名)からメールを検索，返却
//見つからなかった時の処理が必要
func YahooMailRecieve(to string, password string, sub string) ([]byte, error) {

	//urlTarget := "http://***.***.***.***:80/recieve"
	urlTarget := "http://localhost:8080/recieve"
	args := url.Values{}
	args.Add("to", to)
	args.Add("password", password)
	args.Add("subject", sub)

	// Debug
	//fmt.Println("args:\n" + m.From + "\n" + m.To + "\n" + m.Password + "\n" + m.Sub + "\n" + string(m.Msg) + "\n")

	res, err := http.PostForm(urlTarget, args)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, err
	}

	//checkResponse(body)

	return body, nil
}

func checkResponse(body []byte) {
	strJSON := string(body)
	fmt.Println("checkResponse\n" + strJSON)
}

// RecieveMailHandle is recieve mail from POP3 server to web server
func RecieveMailHandle(w http.ResponseWriter, r *http.Request) {

	const yahooPOPServer = "pop.mail.yahoo.co.jp:995"
	const timeout = time.Second * 5

	r.ParseForm()

	postForm := r.PostForm

	to := postForm.Get("to")
	subject := postForm.Get("subject")
	password := postForm.Get("password")

	host, _, _ := net.SplitHostPort(yahooPOPServer)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := pop3.Dial(yahooPOPServer, pop3.UseTLS(tlsconfig), pop3.UseTimeout(timeout))
	if err != nil {
		log.Panic(err)
	}

	if err = conn.Auth(to, password); err != nil {
		log.Panic(err)
	}

	var m *coreMail.MailStruct

	//メールサーバから該当するメールを取り出す
	for i := 1; ; i++ {
		fmt.Println("message number: ", i)

		//計測Retr()
		startRetr := time.Now()
		text, err := conn.Retr(uint32(i))
		endRetr := time.Now()
		fmt.Printf("Retr(): %f\n", (endRetr.Sub(startRetr)).Seconds())
		if err != nil {
			log.Panic(err)
		}
		//fmt.Println("server.go: " + text)

		//計測parseHeader
		startParseHeader := time.Now()
		m = parseHeader(text)
		endParseHeader := time.Now()
		fmt.Printf("parseHeader(): %f\n", (endParseHeader.Sub(startParseHeader)).Seconds())

		//fmt.Println(m.From, "\n", m.To, "\n", m.Sub)
		//fmt.Println(hex.Dump(m.Msg))

		if m.Sub != subject {
			continue
		}

		if err = conn.Dele(uint32(i)); err != nil {
			log.Panic(err)
		}
		if err = conn.Quit(); err != nil {
			log.Panic(err)
		}
		break
	}
	startWrite := time.Now()
	w.Write([]byte(m.Msg))
	endWrite := time.Now()
	fmt.Printf("Write(): %f\n", (endWrite.Sub(startWrite)).Seconds())
}

func parseHeader(text string) *coreMail.MailStruct {
	r := strings.NewReader(text)
	mm, err := mail.ReadMessage(r)
	if err != nil {
		log.Println(err)
	}
	header := mm.Header
	sub := header.Get("Subject")
	//fmt.Println(sub)
	msg, err := ioutil.ReadAll(mm.Body)
	decodedMsg, err := base64.StdEncoding.DecodeString(string(msg))
	//fmt.Println(hex.Dump(msg))
	if err != nil {
		log.Println(err)
	}

	m := coreMail.NewMailStruct("", "", "", "", sub, decodedMsg)

	return m
}
