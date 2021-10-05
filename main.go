package main

import (
  "log"
  "net/http"
)

func main() {
  http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Serving %s\n", r.RemoteAddr)
    w.Header().Add("Content-type", "text/plain")
    w.Write([]byte("Hello World!\n"))
  }))
}
