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

<<<<<<< HEAD
	// Listen start 8080/80 port.
	//8080 local
	//80
	//http.ListenAndServe(":8080", http.DefaultServeMux)
=======
	// Listen start 8080 port.
	http.ListenAndServe(":8080", http.DefaultServeMux)
>>>>>>> 113df24cf0935105cafe609375a8f67144c60d34
	//http.ListenAndServe(":80", http.DefaultServeMux)
}
