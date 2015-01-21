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
		if erro != nil {
			log.Fatalf("error establishing connection, %s", erro)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	//Handle - ini
	defer conn.Close()
	cli_input := bufio.NewReader(conn)
	req, err := http.ReadRequest(cli_input)

	if err != nil {
		log.Printf("error reaading request from client %v \n", err)
		return
	}

	endpoint := strings.Join([]string{req.URL.Host, ":80"}, "")

	backend, err := net.Dial("tcp", endpoint)

	if err != nil {
		log.Printf("error dialing to backend %v \n", err)
		return
	}

	defer backend.Close()

	err = req.Write(backend)

	if err != nil {
		log.Printf("error sending request to backend %v \n", err)
		return
	}

	be_input := bufio.NewReader(backend)
	resp, err := http.ReadResponse(be_input, req)

	if err != nil {
		log.Printf("error reading response from %v \n", err)
		return
	}

	resp.Close = true
	err = resp.Write(conn)

	if err != nil {
		log.Printf("error sending response to client %v \n", err)
		return
	}

	log.Printf("HTTP status %v \n", resp.Status)

	//Handle - Fim
}

func dialRemote(req *http.Request) (net.Conn, error) {
	endPoint := strings.Join([]string{req.URL.Host, ":80"}, "")
	log.Printf("URL - %v \n", endPoint)
	return net.Dial("tcp", endPoint)
}
