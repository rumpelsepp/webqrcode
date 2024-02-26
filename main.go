package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type nopCloser struct {
	io.ReadWriter
}

func (*nopCloser) Close() error {
	return nil
}

func createQR(data string) ([]byte, error) {
	qrc, err := qrcode.New(data)
	if err != nil {
		return nil, err
	}

	var (
		buf     bytes.Buffer
		options = standard.WithBuiltinImageEncoder(standard.PNG_FORMAT)
		w       = standard.NewWithWriter(&nopCloser{&buf}, options)
	)

	if err := qrc.Save(w); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type qrserver struct{}

func (s *qrserver) handleGET(w http.ResponseWriter, r *http.Request) {
	prefix := os.Getenv("WEBQRCODE_PREFIX")

	fmt.Fprintf(w, `
<!DOCTYPE html>
<head>
<title>Steff's QR Code Generator</title>
<meta charset="utf-8">
</head>
<body>
 <form action="%s/qr" method="post">
  <label for="data">Data</label><br>
  <input type="text" id="data" name="data"><br>
  <input type="submit" value="Submit">
</form> 
</body>
`, prefix)
}

func (s *qrserver) handlePOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data := r.FormValue("data")
	qrData, err := createQR(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qrData)
}

func (s *qrserver) serve(listen string) error {
	prefix := os.Getenv("WEBQRCODE_PREFIX")

	r := mux.NewRouter()
	if prefix != "" {
		r = r.PathPrefix(prefix).Subrouter()
	}

	r.HandleFunc("/", s.handleGET)
	r.HandleFunc("/qr", s.handlePOST)

	srv := &http.Server{
		Addr:         listen,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func main() {
	s := qrserver{}
	if err := s.serve(os.Getenv("WEBQRCODE_LISTEN")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
