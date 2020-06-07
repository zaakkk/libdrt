package send

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"../coreMail"
)

//Send Mail parameter to Web server.
func YahooMailSend(m coreMail.MailStruct) error {

	fmt.Printf("*** 開始 ***\n")
	url_target := "http://localhost:8080/send"
	args := url.Values{}
	args.Add("from", m.From)
	args.Add("to", m.To)
	args.Add("password", m.Password)
	args.Add("subject", m.Sub)
	args.Add("body", string(m.Msg))

	// Debug
	//fmt.Println("args:\n" + m.From + "\n" + m.To + "\n" + m.Password + "\n" + m.Sub + "\n" + string(m.Msg) + "\n")

	res, err := http.PostForm(url_target, args)
	if err != nil {
		fmt.Println("Request error:", err)
		//os.Exit(0)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		//os.Exit(0)
		return err
	}

	str_json := string(body)
	fmt.Println("str_json: " + str_json)

	fmt.Printf("*** 終了 ***\n")

	return nil

}
