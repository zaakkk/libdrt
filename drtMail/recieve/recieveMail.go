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

	"github.com/simia-tech/go-pop3"
)

//const yahooPOPServer = "pop.mail.yahoo.co.jp:995"
//const timeout = time.Second * 5

//sub(件名)からメールを検索，返却
//見つからなかった時の処理が必要
func YahooMailRecieve(to string, password string, sub string) ([]byte, error) {

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
		text, err := conn.Retr(uint32(i))
		if err != nil {
			log.Panic(err)
		}
		//fmt.Println("server.go: " + text)
		m = parseHeader(text)

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
	w.Write([]byte(m.Msg))
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
