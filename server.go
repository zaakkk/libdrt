package main

import (
	"net/http"

	"./drtMail/recieve"
	"./drtMail/send"
)

func main() {

	// Open HTML.
	http.Handle("/", http.FileServer(http.Dir(".")))

	// send mail to SMTP server.
	http.HandleFunc("/send", send.SendMailHandle)

	// Recieve mail to POP3 server.
	http.HandleFunc("/recieve", recieve.RecieveMailHandle)

	// Listen start 8080 port.
	http.ListenAndServe(":8080", http.DefaultServeMux)
	//http.ListenAndServe(":80", http.DefaultServeMux)
}
