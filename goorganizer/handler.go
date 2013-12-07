package goorganizer

import (
    "fmt"
    "net/http"
    "time"
)


func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, time.Now())
    fmt.Fprint(w, "Hello, world!")
}

