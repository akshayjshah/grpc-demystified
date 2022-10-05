package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// restHandler takes a RESTful HTTP approach to creating a Pet. It's designed
// to fit on a slide, so it ignores many errors.
func restHandler(w http.ResponseWriter, r *http.Request) {
	const ctype = "application/json"
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
	var pet Pet
	json.NewDecoder(r.Body).Decode(&pet)
	// ✨ save to imaginary DB ✨
	w.Header().Set("Content-Type", ctype)
	json.NewEncoder(w).Encode(pet)
}

// callREST calls the RESTful handler and unmarshals the response.
func callREST(client *http.Client, url string, logger *log.Logger) error {
	body, err := json.Marshal(&Pet{Name: "Fido"})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %v", res.Status)
	}
	if ctype := res.Header.Get("Content-Type"); ctype != "application/json" {
		return fmt.Errorf("unexpected content-type %q", ctype)
	}
	var pet Pet
	if err := json.NewDecoder(res.Body).Decode(&pet); err != nil {
		return err
	}
	logger.Printf("REST response: %+v", pet)
	return nil
}
