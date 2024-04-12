package handler

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type LogRecord struct {
	http.ResponseWriter
	status int
}

func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (w *LogRecord) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		record := &LogRecord{
			ResponseWriter: w,
			status:         200,
		}
		next.ServeHTTP(record, r)
		elapsed := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.RequestURI, record.status, elapsed)
	})
}
