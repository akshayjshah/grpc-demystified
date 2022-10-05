package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// prefix is a gRPC-style envelope. It's designed to fit on a slide, so it
// ignores the first byte of bitwise flags and many errors.
type prefix [5]byte

func (p prefix) Size() int {
	return int(binary.BigEndian.Uint32(p[1:5]))
}

func (p *prefix) SetSize(n int) {
	binary.BigEndian.PutUint32(p[1:5], uint32(n))
}

// grpcHandler implements a gRPC API to create a Pet. It's designed to fit on a
// slide, so it ignores many errors.
func grpcHandler(w http.ResponseWriter, r *http.Request) {
	const ctype = "application/grpc+json"
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != ctype {
		w.Header().Set("Accept", ctype)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	var pre prefix
	r.Body.Read(pre[:])
	input := make([]byte, pre.Size())
	r.Body.Read(input)
	var pet Pet
	json.Unmarshal(input, &pet)
	// ✨ save to imaginary DB ✨
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Te", "trailers") // flush out incompatible proxies
	out, _ := json.Marshal(pet)
	pre.SetSize(len(out))
	w.Write(pre[:])
	w.Write(out)
	w.Header().Set(http.TrailerPrefix+"Grpc-Status", "0")
	w.Header().Set(http.TrailerPrefix+"Grpc-Message", "")
}

// callGRPC calls the gRPC handler and unmarshals the response.
func callGRPC(client *http.Client, baseURL string, logger *log.Logger) error {
	msg, err := json.Marshal(&Pet{Name: "Fido"})
	if err != nil {
		return err
	}
	var pre prefix
	pre.SetSize(len(msg))
	body := append(pre[:], msg...)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/pet.v1.PetService/Create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/grpc+json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %v", res.Status)
	}
	if ctype := res.Header.Get("Content-Type"); ctype != "application/grpc+json" {
		return fmt.Errorf("unexpected content-type %q", ctype)
	}
	pre = prefix{}
	if _, err := res.Body.Read(pre[:]); err != nil {
		return err
	}
	data := make([]byte, pre.Size())
	if _, err := res.Body.Read(data); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	io.Copy(io.Discard, res.Body) // read to EOF to access trailers
	if status := res.Trailer.Get("Grpc-Status"); status != "0" {
		return fmt.Errorf("unexpected gRPC status %q", status)
	}
	var pet Pet
	if err := json.Unmarshal(data, &pet); err != nil {
		return err
	}
	logger.Printf("gRPC response: %+v", pet)
	return nil
}

// jsonCodec is a gRPC codec backed by the standard library's encoding/json.
type jsonCodec struct{}

func (c *jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (c *jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (c *jsonCodec) Name() string {
	return "json"
}

// verifyGRPC calls the gRPC handler using grpc-go, verifying that we've
// implemented the basics of the protocol correctly.
func verifyGRPC(addr string, logger *log.Logger) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	var pet Pet
	if err := conn.Invoke(
		context.Background(),
		"/pet.v1.PetService/Create",
		&Pet{Name: "Fido"},
		&pet,
		grpc.ForceCodecCallOption{Codec: &jsonCodec{}},
	); err != nil {
		return err
	}
	logger.Printf("grpc-go response: %+v", pet)
	return nil
}
