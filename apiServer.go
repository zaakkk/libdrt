package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/zaakkk/libdrt/drtMail/coreMail"
	"github.com/zaakkk/libdrt/drtMail/send"
)

func main() {
	var mailQueue chan coreMail.MailStruct = make(chan coreMail.MailStruct, 300)
	go SendMainLoop(mailQueue)
	Listen(50030, "localhost", mailQueue)
}

func SendMainLoop(queue chan coreMail.MailStruct) {
	for {
		m, ok := <-queue
		if !ok {
			fmt.Println("m is closed")
			return
		}
		fmt.Println("m.Sub: " + m.Sub)
		send.SendMailToSMTP(m)
	}
}

func Listen(port int, host string, queue chan coreMail.MailStruct) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	fmt.Println("waiting for the connection ...")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

ListenLoop:
	for {
		cconn := make(chan *net.TCPConn)
		cerr := make(chan error)
		go func() {
			conn, err := listener.AcceptTCP()
			if err != nil {
				cerr <- err
				return
			}
			cconn <- conn
		}()

		//debug
		//fmt.Println("len(mailQueue): ", len(mailQueue))

		select {
		case conn, ok := <-cconn:
			if !ok {
				fmt.Println("conn is closed")
			}
			go InsertToQueue(conn, queue)
		case err, ok := <-cerr:
			if !ok {
				fmt.Println("err is closed")
			}
			fmt.Println(err)
		case <-interrupt:
			break ListenLoop
		}
		cconn = nil
		cerr = nil
	}
	return nil
}

//apiサーバが受け取ったメールデータを構造体に直してキューに入れる
func InsertToQueue(conn *net.TCPConn, queue chan coreMail.MailStruct) {
	defer conn.Close()

	buf := make([]byte, 1024)
	jsonData := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		jsonData += string(buf[:n])
		if n != 1024 {
			break
		}
	}
	//debug
	//fmt.Println(jsonData)
	var m coreMail.MailStruct
	err := json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		fmt.Println(err)
	}
	//debug
	//fmt.Printf("From: %v, To: %v, Sub: %v\n", m.From, m.To, m.Sub)

	queue <- m
	//debug
	fmt.Println("len(mailQueue): ", len(queue))
}
