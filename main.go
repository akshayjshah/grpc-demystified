package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/v1/pets", http.HandlerFunc(restHandler))
	mux.Handle("/pet.v1.PetService/Create", http.HandlerFunc(grpcHandler))

	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	srv := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}), // support HTTP/2 without TLS
	}
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Fatalf("error listening on localhost: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("error serving HTTP: %v", err)
		}
	}()

	client := &http.Client{Transport: &http.Transport{}}
	if err := callREST(client, "http://localhost:8080/v1/pets", logger); err != nil {
		logger.Fatalf("error calling REST handler: %v", err)
	}
	// gRPC requires HTTP/2, but x/net/http2 makes h2c annoying.
	client.Transport = &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	if err := callGRPC(client, "http://localhost:8080", logger); err != nil {
		logger.Fatalf("error calling RPC handler with hand-written client: %v", err)
	}
	if err := verifyGRPC("localhost:8080", logger); err != nil {
		logger.Fatalf("error calling RPC handler with grpc-go client: %v", err)
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
	wg.Wait()
}
