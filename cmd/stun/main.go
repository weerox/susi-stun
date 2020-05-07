package main

import (
	"crypto/tls"
	"net"
)

const certFile string = "cert.pem"
const keyFile string = "key.pem"

func main() {
	/* run on TCP and UDP on port 3478 and TLS-over-TCP on port 5349 */

	go func() {
		l, err := net.Listen("tcp", ":3478")
		if err != nil {
			// TODO: handle error
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				// TODO: handle error
			}

			// handle the request
			// TODO: move to separate function
			go func(conn net.Conn) {
				defer conn.Close()
			}(conn)
		}
	}()

	go func() {
		pc, err := net.ListenPacket("udp", ":3478")
		if err != nil {
			// TODO: handle error
		}
		defer pc.Close()

		// handle the request
		// TODO: move to separate function
		go func(conn net.PacketConn) {
			defer conn.Close()
		}(pc)
	}()

	go func() {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			// TODO: handle error
		}

		conf := &tls.Config{Certificates: []tls.Certificate{cert}}

		l, err := tls.Listen("tcp", ":5349", conf)
		if err != nil {
			// TODO: handle error
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				// TODO: handle error
			}

			// handle request
			// TODO: move to separate function
			go func(conn net.Conn) {
				defer conn.Close()
			}(conn)
		}
	}()

	select {}
}
