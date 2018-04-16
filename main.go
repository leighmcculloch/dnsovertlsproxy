package main // import "4d63.com/dnsovertlsproxy"

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const dnsHost = "1.1.1.1:853"

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

	for {
		query := make([]byte, 1232)

		n, addr, err := conn.ReadFromUDP(query)
		if err != nil {
			log.Println("error", err)
			continue
		}
		log.Println("received query from", addr)
		query = query[:n]

		go func() {
			resp, err := dns(query)
			if err != nil {
				log.Println("error", err)
				return
			}

			_, err = conn.WriteToUDP(resp, addr)
			if err != nil {
				log.Println("error", err)
				return
			}
			log.Println("sent results to", addr)
		}()
	}
}

func dns(query []byte) ([]byte, error) {
	conn, err := tls.Dial("tcp", dnsHost, &tls.Config{})
	if err != nil {
		return nil, err
	}

	req := make([]byte, len(query)+2)
	binary.BigEndian.PutUint16(req[0:2], uint16(len(query)))
	copy(req[2:], query)

	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}
	log.Println("sent query on to dns host", dnsHost)

	resp, err := dnsReadResponse(conn)
	if err != nil {
		return nil, err
	}
	log.Println("received response from dns host", dnsHost, len(resp), "bytes")

	return resp, nil
}

func dnsReadResponse(r io.Reader) ([]byte, error) {
	length, err := dnsReadLength(r)
	if err != nil {
		return nil, err
	}

	return dnsReadMessage(r, length)
}

func dnsReadLength(r io.Reader) (int, error) {
	bytes := make([]byte, 2)
	n, err := r.Read(bytes)
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, fmt.Errorf("reading length did not receive enough bytes")
	}
	length := int(binary.BigEndian.Uint16(bytes))
	return length, nil
}

func dnsReadMessage(r io.Reader, length int) ([]byte, error) {
	resp := make([]byte, length)
	offset := 0
	for offset < length {
		n, err := r.Read(resp[offset:])
		if err != nil {
			return nil, err
		}
		offset += n
	}
	return resp, nil
}
