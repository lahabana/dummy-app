package main

import (
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()
	defer func() {
		l.Infow("Done shutting down service!")
	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lsn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	l.Infow("Running on", "addr", lsn.Addr())

	err = http.Serve(lsn, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Infow("Got request", "method", r.Method, "path", r.URL.String(), "content-length", r.ContentLength)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Infow("Failed request", err.Error())
			_, _ = w.Write([]byte("Bad error"))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
		}
	}))
	if err != nil {
		panic(err)
	}
}
