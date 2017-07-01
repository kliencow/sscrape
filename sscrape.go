package sscrape

import (
	"fmt"
	"net/http"
	"log"
	"bytes"
	"strconv"
	"net/url"
	"strings"
	"os"
)


type TargetServer struct {
	Jar []*http.Cookie
	Host string
	ConfLoc string
}

func (ts TargetServer) URL(path string) (string, error) {
	u, err := url.ParseRequestURI(ts.Host)
	if err != nil {
		return "", fmt.Errorf("url for host unable to be parsed: %v", err)
	}

	u.Path = path
	return u.String(), nil
}

func (ts TargetServer) PostRequest(path string, form url.Values) (*http.Request, error) {
	url, err := ts.URL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("invalid request, unable to construct: %v", err)
	}

	req.Header.Set("User-Agent", "SScraper/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (ts TargetServer) SwapCookies(foundCookies []*http.Cookie, cookiePrefix string) {
	newJar := make([]*http.Cookie, len(ts.Jar), len(ts.Jar)+1)

	for _, cookie := range(ts.Jar) {
		if strings.HasPrefix(cookie.Name, cookiePrefix) {
			//
		} else {
			newJar = append(newJar, cookie)
		}
	}

	for _, cookie := range(foundCookies) {
		if strings.HasPrefix(cookie.Name, cookiePrefix) {
			newJar = append(newJar, cookie)
		}
	}

	ts.Jar = newJar
}

func (ts TargetServer) HasCookie(cookies []*http.Cookie, cookiePrefix string) bool {
	for _, cookie := range(cookies) {
		if strings.HasPrefix(cookie.Name, cookiePrefix) {
			return true
		}
	}

	return false
}




