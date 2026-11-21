package media_test

import (
	"testing"

	"github.com/fogfish/gouldian/v3/internal/media"
)

func TestIsBinary(t *testing.T) {
	for _, tt := range []struct {
		ct     string
		body   []byte
		binary bool
	}{
		{"", nil, false},
		{"", []byte("text"), false},
		{"", []byte{0xff, 0xfe}, true},
		{"text/plain", []byte("text"), false},
		{"text/plain; charset=utf-8", []byte("text"), false},
		{"text/plain", []byte{0xff}, true},
		{"application/json", []byte(`{}`), false},
		{"application/problem+json", []byte(`{}`), false},
		{"application/vnd.api+json", []byte(`{}`), false},
		{"application/octet-stream", []byte("abc"), true},
		{"image/png", []byte{0x89, 'P'}, true},
		{"application/pdf", []byte("%PDF"), true},
		{"\x00bad", []byte("abc"), false},
	} {
		if got := media.IsBinary(tt.ct, tt.body); got != tt.binary {
			t.Errorf("IsBinary(%q, %q) = %v, want %v", tt.ct, tt.body, got, tt.binary)
		}
	}
}
