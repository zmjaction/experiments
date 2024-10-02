package main

import (
	"fmt"
	"net/http"
	"strings"
)

// http request header
// Cookie : session_id=123; value=hello-world;lang=en;lang=zh-CN
func parseCookie() {
	lines := "session_id=123; value=hello-world;lang=en;lang=zh-CN"
	cookies, _ := http.ParseCookie(lines)
	for _, cookie := range cookies {
		fmt.Printf("%s: %s\n", cookie.Name, cookie.Value)
	}
	fmt.Println(strings.Repeat("-", 50))
}

// http response header
// Set-Cookie: session_id=123; MaxAge=0;lang=zh=CN;Domain=.123.com
func paeseSetCookie() {
	line := "session_id=123; MaxAge=0;lang=zh=CN;Domain=.123.com"
	cookie, _ := http.ParseSetCookie(line)
	fmt.Println("Name:", cookie.Name)
	fmt.Println("Name:", cookie.Value)
	fmt.Println("Name:", cookie.Domain)
	fmt.Println(strings.Repeat("-", 50))
}
func main() {
	parseCookie()
	paeseSetCookie()
}
