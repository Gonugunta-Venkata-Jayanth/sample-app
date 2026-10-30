package main

import (
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestColorHandlers(t *testing.T) {
	cases := map[string]struct {
		handler http.HandlerFunc
		r, g, b uint32
	}{
		"/blue": {blueHandler, 0, 0, 0xffff},
		"/red":  {redHandler, 0xffff, 0, 0},
	}
	for path, tc := range cases {
		rec := httptest.NewRecorder()
		tc.handler(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "image/png" {
			t.Fatalf("%s: status %d, content-type %q", path, rec.Code, rec.Header().Get("Content-Type"))
		}
		img, err := png.Decode(rec.Body)
		if err != nil {
			t.Fatalf("%s: invalid png: %v", path, err)
		}
		if b := img.Bounds(); b.Dx() != 100 || b.Dy() != 100 {
			t.Fatalf("%s: unexpected size %v", path, b)
		}
		r, g, b, _ := img.At(50, 50).RGBA()
		if r != tc.r || g != tc.g || b != tc.b {
			t.Fatalf("%s: got colour %d,%d,%d", path, r, g, b)
		}
	}
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("healthz: %d %q", rec.Code, rec.Body.String())
	}
}
