package send

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"../coreMail"
)

//YahooMailSend is send mail parameter to Web server.
func YahooMailSend(m coreMail.MailStruct) error {

	//fmt.Printf("*** 開始 ***\n")

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

	/*
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Request error:", err)
			return err
		}

		strJSON := string(body)
		fmt.Println("strJSON: " + strJSON)
	*/

	//fmt.Printf("*** 終了 ***\n")

	return nil

}
