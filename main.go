package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

var (
	listenAddr = flag.String("l", ":3000", "Listening address")
	secure     = flag.Bool("secure", false, "Serving using HTTPS (self-signed cert)")
)

func genTLSConfig() (*tls.Config, error) {
	keypair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(180 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tpl, tpl, keypair.Public(), keypair)
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{certDER},
				PrivateKey:  keypair,
			},
		},
	}
	return conf, nil
}

func main() {
	flag.Parse()
	var ln net.Listener
	var err error
	if *secure {
		var conf *tls.Config
		if conf, err = genTLSConfig(); err != nil {
			log.Fatalf("Generate TLS config: %v", err)
		}
		ln, err = tls.Listen("tcp", *listenAddr, conf)
	} else {
		ln, err = net.Listen("tcp", *listenAddr)
	}
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Serving %s\n", r.RemoteAddr)
			w.Header().Add("Content-type", "text/plain")
			w.Write([]byte("RQV-r0: Hello World!\n"))
		}),
	}
	server.Serve(ln)
}
