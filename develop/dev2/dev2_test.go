package main

import (
	"log/slog"
	"testing"
)

func TestUnzip(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"a4bc2d5e", "a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd", "abcd"},
		{"45", "45", ""},
		{"qwe\\4\\5", "qwe\\4\\5", "qwe45"},
		{"qwe\\\\5", "qwe\\\\5", "qwe\\\\\\\\\\"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unzip(tt.input)
			if err != nil {
				slog.Error(err.Error())
			}
			if result != tt.want {
				t.Errorf("want: %s, got: %s", tt.want, result)
			}
		})
	}
}
