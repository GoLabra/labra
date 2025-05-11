package strcase_test

import (
	"testing"

	"github.com/GoLabra/labrago/src/api/strcase"
)

func TestToCamel(t *testing.T) {
	var tts = []struct {
		in  string
		out string
	}{
		{"a test", "aTest"},
		{"aTest", "aTest"},
		{"twoTest", "twoTest"},
		{"TwoTest", "TwoTest"},
		{"two Test", "twoTest"},
		{"two test", "twoTest"},
		{"two_test", "twoTest"},
		{"two_Test", "twoTest"},
		{"two_TeSt", "twoTeSt"},
	}

	for _, tt := range tts {
		v := strcase.ToCamel(tt.in)
		if strcase.ToCamel(tt.in) != tt.out {
			t.Errorf("got %q, want %q", v, tt.out)
		}
	}
}
func TestToLowerCamel(t *testing.T) {
	var tts = []struct {
		in  string
		out string
	}{
		{"a test", "aTest"},
		{"aTest", "aTest"},
		{"twoTest", "twoTest"},
		{"TwoTest", "twoTest"},
		{"two Test", "twoTest"},
		{"two test", "twoTest"},
		{"two_test", "twoTest"},
		{"two_Test", "twoTest"},
		{"two_TeSt", "twoTeSt"},
	}

	for _, tt := range tts {
		v := strcase.ToLowerCamel(tt.in)
		if strcase.ToLowerCamel(tt.in) != tt.out {
			t.Errorf("got %q, want %q", v, tt.out)
		}
	}
}

func TestToSnake(t *testing.T) {
	var tts = []struct {
		in  string
		out string
	}{
		{"a test", "a_test"},
		{"aTest", "a_test"},
		{"twoTest", "two_test"},
		{"two Test", "two_test"},
		{"two test", "two_test"},
		{"two_test", "two_test"},
		{"two_Test", "two_test"},
		{"two_TeSt", "two_te_st"},
	}

	for _, tt := range tts {
		v := strcase.ToSnake(tt.in)
		if strcase.ToSnake(tt.in) != tt.out {
			t.Errorf("got %q, want %q", v, tt.out)
		}
	}
}

func TestToPascal(t *testing.T) {
	var tts = []struct {
		in  string
		out string
	}{
		{"a test", "a_test"},
		{"aTest", "a_test"},
	}

	for _, tt := range tts {
		v := strcase.ToSnake(tt.in)
		if strcase.ToSnake(tt.in) != tt.out {
			t.Errorf("got %q, want %q", v, tt.out)
		}
	}
}
