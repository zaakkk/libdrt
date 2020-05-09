package recieve

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"strings"
	"time"

	"../../go-pop3"
	"../coreMail"
)

const yahooPOPServer = "pop.mail.yahoo.co.jp:995"
const timeout = time.Second * 5

//sub(件名)からメールを検索，返却
//見つからなかった時の処理が必要
func YahooMailRecieve(mailAddress string, password string, sub string) *coreMail.MailStruct {
	host, _, _ := net.SplitHostPort(yahooPOPServer)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := pop3.Dial(yahooPOPServer, pop3.UseTLS(tlsconfig), pop3.UseTimeout(timeout))
	if err != nil {
		log.Panic(err)
	}

	if err = conn.Auth(mailAddress, password); err != nil {
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

		m = parseHeader(text)

		//fmt.Println(m.From, "\n", m.To, "\n", m.Sub)
		//fmt.Println(hex.Dump(m.Msg))

		if m.Sub != sub {
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

	return m
}

//取り出したメールから更に本文とタイトルを取り出す
func parseHeader(text string) *coreMail.MailStruct {
	r := strings.NewReader(text)
	mm, err := mail.ReadMessage(r)
	if err != nil {
		log.Println(err)
	}
	header := mm.Header
	sub := header.Get("Subject")
	msg, err := ioutil.ReadAll(mm.Body)
	decodedMsg, err := base64.StdEncoding.DecodeString(string(msg))
	//fmt.Println(hex.Dump(msg))
	if err != nil {
		log.Println(err)
	}

	//m := coreMail.NewMailStruct("", "", "", "", sub, msg)
	m := coreMail.NewMailStruct("", "", "", "", sub, decodedMsg)

	return m
}
