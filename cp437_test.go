package xploader

import "testing"

func TestCP437Decoder(t *testing.T) {
	if got := CP437Decoder(0); got != '\u0000' {
		t.Errorf("Expected Null for code 0, got %q", got)
	}
	if got := CP437Decoder(255); got != '□' {
		t.Errorf("Expected □ for code 255, got %q", got)
	}
}

func TestCP437Encoder(t *testing.T) {
	if code := CP437Encoder('\u0000'); code != 0 {
		t.Errorf("Expected code 0 for Null, got %d", code)
	}
	if code := CP437Encoder('\u00A0'); code != 255 {
		t.Errorf("Expected code 255 for NBSP, got %d", code)
	}
}

func TestEncoderFallback(t *testing.T) {
	unmapped := rune(0x10FFFF) // unlikely to exist
	if code := CP437Encoder(unmapped); code != unmapped {
		t.Errorf("Expected fallback for unmapped rune, got %d", code)
	}
}

func TestCP437DecoderFallback(t *testing.T) {
	// Pick a value outside the CP437 range that is not in the map
	code := int32(999)
	got := CP437Decoder(code)
	if got != rune(code) {
		t.Errorf("Expected fallback to return %d as rune, got %q", code, got)
	}
}
