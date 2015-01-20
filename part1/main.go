package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("erro creating server, %s", err)
	}

	for {
		conn, erro := ln.Accept()
		if err != nil {
			log.Fatalf("error establishing connection, %s", erro)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, _ := http.ReadRequest(bufio.NewReader(conn))
	log.Printf("printing the request %v \n", req)
	rConn, err := dialRemote(req)

	if err == nil {
		defer rConn.Close()
		res := new(http.Response)
		log.Println("reading response")
		res, err = http.ReadResponse(bufio.NewReader(rConn), req)
		defer res.Body.Close()

		log.Printf("printing the RESPONSE %v \n", res)
		res.Write(bufio.NewWriter(conn))
	} else {
		log.Printf("error sending request to backend %v \n", err)
	}
}

func dialRemote(req *http.Request) (net.Conn, error) {
	endPoint := strings.Join([]string{req.URL.Host, ":80"}, "")
	log.Printf("URL - %v \n", endPoint)
	return net.Dial("tcp", endPoint)
}
