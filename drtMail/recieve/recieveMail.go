package recieve

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

//sub(件名)からメールを検索，返却
//見つからなかった時の処理が必要
func YahooMailRecieve(to string, password string, sub string) ([]byte, error) {

	urlTarget := "http://13.231.164.236:80/recieve"
	//urlTarget := "http://localhost:8080/recieve"
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

	decodedMsg, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println(err)
	}

	return decodedMsg, nil
}

func checkResponse(body []byte) {
	strJSON := string(body)
	fmt.Println("checkResponse\n" + strJSON)
}

func RecieveMailHandle(w http.ResponseWriter, r *http.Request) {

	const googleImapServer = "imap.gmail.com:993"

	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	postForm := r.PostForm

	username := postForm.Get("to")
	targetSubject := postForm.Get("subject")
	password := postForm.Get("password")

	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(googleImapServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the all message
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	//seqSet.AddNum(mbox.Messages)
	seqSet.AddRange(from, to)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {

		//msg := <-messages
		if msg == nil {
			log.Fatal("Server didn't returned message")
		}

		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		// Print some info about the message
		header := mr.Header
		subject := ""
		if subject, err = header.Subject(); err == nil {
			log.Println("Subject:", subject)
		}

		if targetSubject != subject {
			continue
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				log.Println("Got text: %v", string(b))
				w.Write(b)
			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				log.Println("Got attachment: %v", filename)
			}
		}
	}
}
