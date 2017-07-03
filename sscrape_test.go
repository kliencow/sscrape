package sscrape

import (
	"testing"
	"net/url"
	"io/ioutil"
)

func TestURL(t *testing.T) {
	ts := TargetServer{
		Host: "http://example.com",
	}

	u, err := ts.URL("path/to/file.txt")
	if err != nil {
		t.Errorf("unexpected error building url: %v", err)
	}

	eu := "http://example.com/path/to/file.txt"
	if u != eu {
		t.Errorf("failed to build correct url expected %v got %v.", u, eu)
	}
}

func TestBasicGetRequest(t *testing.T) {
	ts := TargetServer{
		Host: "http://example.com",
	}

	req, err := ts.Request("GET", "form.php", nil)
	if err != nil {
		t.Errorf("unable to build legit request, %v", err)
	}

	if len(req.Cookies()) > 0 {
		t.Errorf("Cookies set on request with no cookies in the jar")
	}

	if req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		t.Errorf("Content type not set correctly for GET")
	}
}

func TestParamsPostRequest(t *testing.T) {
	ts := TargetServer{
		Host: "http://example.com",
	}

	form := url.Values{
		"a":{"foo"},
		"b":{"bar"},
	}

	req, err := ts.Request("POST", "form.php", form)
	if err != nil {
		t.Errorf("unable to build legit request, %v", err)
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Content type not set correctly for POST")
	}

	uriStr := req.URL.String()
	expectedUriStr := "http://example.com/form.php"
	if uriStr != expectedUriStr {
		t.Errorf("get uri incorrect. expected '%v' got '%v", expectedUriStr, uriStr)
	}

	expectedBody := "a=foo&b=bar"
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("Unable to parse body")
	}
	if string(body) != expectedBody {
		t.Errorf("Post body incorrect. expected '%v' got '%v", expectedBody, string(body))
	}

}

func TestParamsGetRequest(t *testing.T) {
	ts := TargetServer{
		Host: "http://example.com",
	}

	form := url.Values{
		"a":{"foo"},
		"b":{"bar"},
	}

	req, err := ts.Request("GET", "form.php", form)
	if err != nil {
		t.Errorf("unable to build legit request, %v", err)
	}

	expectedBody := ""
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("unable to parse body")
	}

	uriStr := req.URL.String()
	expectedUriStr := "http://example.com/form.php?a=foo&b=bar"
	if uriStr != expectedUriStr {
		t.Errorf("get uri incorrect. expected '%v' got '%v", expectedUriStr, uriStr)
	}

	if string(body) != expectedBody {
		t.Errorf("get body incorrect. expected '%v' got '%v", expectedBody, string(body))
	}
}