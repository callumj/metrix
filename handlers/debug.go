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

func TestHandler(c http.ResponseWriter, req *http.Request) {
	body := `<h2>Hello, you've been metrixed</h2>
  <img src="/metric/increment?key=test_image&app=metrix&image=yes" />
  <img src="/metric/increment?key=test_image_redirect&app=metrix&redirect=https://graph.facebook.com/callumj/picture?type=large" />`
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", strconv.Itoa(len(body)))
	io.WriteString(c, body)
}
