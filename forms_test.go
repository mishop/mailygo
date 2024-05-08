package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func Test_sanitizeForm(t *testing.T) {
	t.Run("Sanitize form", func(t *testing.T) {
		result := sanitizeForm(&url.Values{"<b>Test</b>": {"<a href=\"https://example.com\">Test</a>"}})
		want := FormValues{"Test": {"Test"}}
		if !reflect.DeepEqual(*result, want) {
			t.Error()
		}
	})
}

func TestFormHandler(t *testing.T) {
	t.Run("GET request to FormHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()
		FormHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Error()
		}
	})
	t.Run("POST request to FormHandler", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://example.com/", nil)
		w := httptest.NewRecorder()
		FormHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusCreated {
			t.Error()
		}
	})
	t.Run("Wrong method request to FormHandler", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://example.com/", nil)
		w := httptest.NewRecorder()
		FormHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Error()
		}
	})
}

func Test_isBot(t *testing.T) {
	t.Run("No bot", func(t *testing.T) {
		os.Clearenv()
		formData := FormValues{"_t_email": {""}, "_replyTo": {"misho@misho.net"}}
		result := isBot(&formData)
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("No honeypot", func(t *testing.T) {
		os.Clearenv()
		formData := FormValues{}
		result := isBot(&formData)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("Bot with empty _replyTo", func(t *testing.T) {
		os.Clearenv()
		formData := FormValues{"_t_email": {""}, "_replyTo": {""}}
		result := isBot(&formData)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})
}

func Test_sendResponse(t *testing.T) {
	t.Run("No redirect", func(t *testing.T) {
		values := &FormValues{}
		w := httptest.NewRecorder()
		sendResponse(values, w)
		if w.Code != http.StatusCreated {
			t.Error()
		}
	})
	t.Run("No redirect 2", func(t *testing.T) {
		values := &FormValues{
			"_redirectTo": {""},
		}
		w := httptest.NewRecorder()
		sendResponse(values, w)
		if w.Code != http.StatusCreated {
			t.Error()
		}
	})
	t.Run("No redirect 3", func(t *testing.T) {
		values := &FormValues{
			"_redirectTo": {"abc", "def"},
		}
		w := httptest.NewRecorder()
		sendResponse(values, w)
		if w.Code != http.StatusCreated {
			t.Error()
		}
	})
	t.Run("Redirect", func(t *testing.T) {
		values := &FormValues{
			"_redirectTo": {"https://example.com"},
		}
		w := httptest.NewRecorder()
		sendResponse(values, w)
		if w.Code != http.StatusSeeOther {
			t.Error()
		}
		if w.Header().Get("Location") != "https://example.com" {
			t.Error()
		}
	})
}
