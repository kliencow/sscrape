// Copyright 2017 Walt Norblad. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sscrape

import (
	"fmt"
	"net/http"
	"bytes"
	"net/url"
	"strings"
	"io/ioutil"
)

// TargetServer represents the server to be scraped
type TargetServer struct {
	// Jar (optional) is the cookie jar used to hold session cookies and other cookies. Normally this is filled
	// with the LoginFunction. It's only really used for cookie based session tracking.
	Jar []*http.Cookie

	// AgentName (optional) is a fun little morsel meant just to give your script a name while it scrapes
	// sites. It's only really used for fun, though I guess you could masquerade as real browsers
	// with it. But why?
	AgentName string

	// Host (required) is the fully qualified URL of the server to be scraped
	// For example: http://example.com:8080 or http://google.com
	Host string

	// SessionCookieName (optional) is the name of the session cookie the target server uses.
	// Providing this cookie name will make Form based logins easier as it will error if the
	// appropriate cookie is not found. Otherwise, it's not used.
	SessionCookieName string
}

// FormLogin tries to log into the site using the login parameters provided in the form parameter.
// Any cookies returned will be stored in the Jar. If a SessionCookieName value is set, this function
// will look for it in the response in particular. If it's not found then an error is returned. This
// is useful for determining if a login was actually successful
func (ts *TargetServer) FormLogin(loginPath string, form url.Values) error {
	req, err := ts.Request("POST", loginPath, form)
	if err != nil {
		return fmt.Errorf("unable to create log in request, %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request to target host due to %v", err)
	}

	if ts.SessionCookieName != "" && ts.HasCookie(resp.Cookies(), ts.SessionCookieName) {
		ts.Jar = resp.Cookies()
	} else {
		return fmt.Errorf("session cookie not found after login due to possible bad username and password")
	}

	return nil
}


// GetPage gets a page from the server and returns the result as a string. At the moment
// there is no support for binary results.
func (ts *TargetServer) GetPage(path string, query url.Values) (string, error) {
	req, err := ts.Request("GET", path, query)
	if err != nil {
		return "", fmt.Errorf("unable to create log in request, %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to make request to target host due to %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body from response, %v", err)
	}

	return  string(body), nil
}


// URL is a helper function that just builds the right urls for me because using net/url
// is tedious.
func (ts TargetServer) URL(path string) (string, error) {
	u, err := url.ParseRequestURI(ts.Host)
	if err != nil {
		return "", fmt.Errorf("url for host unable to be parsed: %v", err)
	}

	u.Path = path
	return u.String(), nil
}

// Reqest build a request and adds common thing that I always end up adding. More importantly
// this loads the cookies into the request
func (ts TargetServer) Request(method string, path string, form url.Values) (*http.Request, error) {
	url, err := ts.URL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("invalid request, unable to construct: %v", err)
	}

	// While this function only accepts text reponses, I don't want to list them at the moment.
	// one day when I feel keep on it, I'll update this with actual accept types.
	req.Header.Set("Accept", "*/*")

	if ts.AgentName == "" {
		req.Header.Set("User-Agent", "SScraper/1.0")
	} else {
		req.Header.Set("User-Agent", ts.AgentName)
	}

	if method != "GET" || method != "HEAD" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for _, cookie := range(ts.Jar) {
		req.AddCookie(cookie)
	}

	return req, nil
}

// SwapsCookies swaps a cookie in a list of cookies. This was supposed to be used to refresh the
// Session cookie if it were dynamic, but then it turned out to just be easier to dump all of the
// cookies in to the jar. One day, I plan on getting back into using this and dealing with cookies
// with intelligence rather than a caveman club.
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

// HasCookie checks a list of cookies if a cooke is there by name.
func (ts TargetServer) HasCookie(cookies []*http.Cookie, cookiePrefix string) bool {
	for _, cookie := range(cookies) {
		if strings.HasPrefix(cookie.Name, cookiePrefix) {
			return true
		}
	}

	return false
}
