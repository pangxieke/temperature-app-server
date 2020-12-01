package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type response struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (r *response) Write(data []byte) (int, error) {
	r.Body = data
	return r.ResponseWriter.Write(data)
}

func (r *response) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now().UnixNano()

		log.Printf("Started %s %s %s", r.Method, r.URL.Path, r.URL.RawQuery)
		headers := []string{}
		for k, v := range r.Header {
			headers = append(headers, fmt.Sprintf("%v: %v", k, v[0]))
		}
		log.Printf("resq headers: %s", strings.Join(headers, ", "))
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		log.Printf("resq body: %s", body)

		// Call the next handler in the chain.
		res := response{ResponseWriter: w}
		next.ServeHTTP(&res, r)

		elapse := float64(time.Now().UnixNano()-begin) / 1000000.0
		headers = []string{}
		for k, v := range res.Header() {
			headers = append(headers, fmt.Sprintf("%v: %v", k, v[0]))
		}
		log.Printf("resp headers: %s", strings.Join(headers, ", "))
		log.Printf("resp: %s\n", res.Body)
		log.Printf("Completed %s %s %s with %d, in %.2f ms \n\n", r.Method, r.URL.Path, r.URL.RawQuery, res.Status, elapse)
	})
}
