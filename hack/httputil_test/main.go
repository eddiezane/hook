package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var foo = `POST / HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 7
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.66.0

foo=bar
`

func main() {
	reader := strings.NewReader(foo)
	r, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		panic(err)
	}
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(string(s))
}
