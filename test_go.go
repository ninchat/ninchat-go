// +build test,!js

package ninchat

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"reflect"
	"unsafe"
)

func init() {
	if s := os.Getenv("NINCHAT_TEST_ADDRESS"); s != "" {
		defaultAddress = s
	}

	if certFile := os.Getenv("NINCHAT_TEST_CLIENT_CERT"); certFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, os.Getenv("NINCHAT_TEST_CLIENT_KEY"))
		if err != nil {
			panic(err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	if os.Getenv("NINCHAT_TEST_INSECURE") == "1" {
		tlsConfig.InsecureSkipVerify = true
	}
}

func (s *Session) testDisconnect() (ok bool) {
	tlsConn := s.test.ws.conn.UnderlyingConn().(*tls.Conn)
	fieldAddr := reflect.ValueOf(tlsConn).Elem().FieldByName("conn").UnsafeAddr()
	netConn := *(*net.Conn)(unsafe.Pointer(fieldAddr))
	tcpConn := netConn.(*net.TCPConn)

	// Corrupt the outgoing TLS stream.
	if _, err := tcpConn.Write([]byte{42}); err != nil {
		log.Printf("TestDisconnect: %v", err)
		return
	}

	if err := tcpConn.CloseWrite(); err != nil {
		log.Printf("TestDisconnect: %v", err)
		return
	}

	ok = true
	return
}
