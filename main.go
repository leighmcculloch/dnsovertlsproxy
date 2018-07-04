package main // import "4d63.com/dnsovertlsproxy"

import (
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var version = "<not set>"

const defaultServerAddr = "1.1.1.1:853"
const defaultListenAddr = ":53"

func main() {
	flagListenAddr := flag.String("listen", defaultListenAddr, "")
	flagServerAddr := flag.String("server", defaultServerAddr, "")
	flagPrintHelp := flag.Bool("help", false, "print this help")
	flagPrintVersion := flag.Bool("version", false, "print version")
	flagVerbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "dnsovertlsproxy is a simple DNS over TLS proxy.\n")
		fmt.Fprintf(os.Stderr, "Usage: dnsovertlsproxy\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagPrintHelp {
		flag.Usage()
		return
	}

	if *flagPrintVersion {
		fmt.Println("dnsovertlsproxy", "v"+version)
		return
	}

	verbose := *flagVerbose
	listenAddr := *flagListenAddr
	serverAddr := *flagServerAddr

	laddr, err := net.ResolveUDPAddr("udp", listenAddr)
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
		if verbose {
			log.Println("received query from", addr)
		}
		query = query[:n]

		resp, err := dns(verbose, serverAddr, query)
		if err != nil {
			log.Println("error", err)
			continue
		}

		_, err = conn.WriteToUDP(resp, addr)
		if err != nil {
			log.Println("error", err)
			continue
		}
		if verbose {
			log.Println("sent results to", addr)
		}
	}
}

func dns(verbose bool, serverAddr string, query []byte) ([]byte, error) {
	conn, err := tls.Dial("tcp", serverAddr, &tls.Config{})
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
	if verbose {
		log.Println("sent query on to server", serverAddr)
	}

	resp, err := dnsReadResponse(conn)
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Println("received response from server", serverAddr, len(resp), "bytes")
	}

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
