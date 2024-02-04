package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKratos(t *testing.T) {
	w := httptest.NewRecorder()

	// Simulation of custom http.DefaultResponseEncoder
	var encode = func(rw http.ResponseWriter, v []byte) {
		// Bad use case. The correct order is 2 -> 1 -> 3
		rw.WriteHeader(500)             // 1
		rw.Header().Set("X-Foo", "foo") // 2
		rw.Write(v)                     // 3
	}
	// Simulation of http.wrapper.Result
	var result = func(code int, v []byte) {
		rw := &responseWriter{w: w}
		rw.WriteHeader(code)
		encode(rw, v)
	}

	result(400, []byte("a"))

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 500 {
		t.Fatal("code not changed")
	}
	if foo := resp.Header.Get("X-Foo"); foo != "" {
		t.Fatal("header should empty when calling WriterHeader with a wrong order")
	}
}

// responseWriter creates a response writer that "caches" the statusCode.
type responseWriter struct {
	code int
	w    http.ResponseWriter
}

func (w *responseWriter) Header() http.Header        { return w.w.Header() }
func (w *responseWriter) WriteHeader(statusCode int) { w.code = statusCode }
func (w *responseWriter) Write(data []byte) (int, error) {
	w.w.WriteHeader(w.code)
	return w.w.Write(data)
}

func TestMy(t *testing.T) {
	w := httptest.NewRecorder()

	// Simulation of custom http.DefaultResponseEncoder
	var encode = func(rw http.ResponseWriter, v []byte) {
		// Bad use case. The correct order is 2 -> 1 -> 3
		rw.WriteHeader(500)             // 1
		rw.Header().Set("X-Foo", "foo") // 2
		rw.Write(v)                     // 3
	}
	// Simulation of http.wrapper.Result
	var result = func(code int, v []byte) {
		// No need to call WriteHeader. Just set the field.
		encode(&myWriter{code, w}, v)
	}

	result(400, []byte("a"))

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 500 {
		t.Fatal("code not changed")
	}
	if foo := resp.Header.Get("X-Foo"); foo != "" {
		t.Fatal("header should empty when calling WriterHeader with a wrong order")
	}
}

// myWriter creates a response writer that changes the default status code.
type myWriter struct {
	code int
	w    http.ResponseWriter
}

func (w *myWriter) Header() http.Header { return w.w.Header() }

func (w *myWriter) WriteHeader(statusCode int) {
	w.code = 0
	w.w.WriteHeader(statusCode)
}
func (w *myWriter) Write(data []byte) (int, error) {
	if w.code != 0 {
		w.w.WriteHeader(w.code)
	}
	return w.w.Write(data)
}
