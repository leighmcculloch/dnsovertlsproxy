package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	laddr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1232)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Print("error ", err)
			continue
		}
		log.Print("received query from ", addr, "[", base64.StdEncoding.EncodeToString(buf[:n]), "]")

		body := bytes.NewReader(buf[:n])
		req, err := http.NewRequest("POST", "https://cloudflare-dns.com/dns-query", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Accept", "application/dns-udpwireformat")
		req.Header.Add("Content-Type", "application/dns-udpwireformat")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Print("error ", err)
			continue
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("error ", err)
			continue
		}

		log.Print("received from 1.1.1.1", "[", base64.StdEncoding.EncodeToString(respBody), "]")
		_, err = conn.WriteToUDP(respBody, addr)
		if err != nil {
			log.Print("error ", err)
			continue
		}
		log.Print("sent results to ", addr)
	}
}
