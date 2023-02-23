package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheushelper "github.com/wizact/go-prometheus/prometheus-helper"
)

func main() {
	port := 8080
	readTimeout := time.Millisecond * 500
	writeTimeout := time.Millisecond * 500

	httpMetrics := prometheushelper.NewHttpMetrics()

	instrumentedHash := prometheushelper.InstrumentHandler(hashFunc, httpMetrics)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/hash", instrumentedHash)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "No tracking for this request")
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	s.ListenAndServe()
}

func hashFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Body != nil {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, "%s", computeSum(body))
	}
}

func computeSum(body []byte) []byte {
	h := sha256.New()
	h.Write(body)
	hashed := hex.EncodeToString(h.Sum(nil))
	return []byte(hashed)
}
