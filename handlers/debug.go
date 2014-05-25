package handlers

import (
	"io"
	"net/http"
	"strconv"
)

func PingHandler(c http.ResponseWriter, req *http.Request) {
	body := "<h2>Pong!</h2>"
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", strconv.Itoa(len(body)))
	io.WriteString(c, body)
}
