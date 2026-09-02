// Package banner отвечает за загрузку и разбор файлов-баннеров (standard.txt,
// shadow.txt, thinkertoy.txt) и предоставляет доступ к ASCII-блоку (8 строк)
// для каждого печатного символа.
//
// Package banner is responsible for loading and parsing banner files
// (standard.txt, shadow.txt, thinkertoy.txt) and provides access to the
// 8-line ASCII block for each printable character.
package banner

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// buildValidBannerFixture builds a well-formed banner file's lines (blank
// separator + Height content lines, repeated Count times) along with the
// Banner Load() should produce from them. Content lines are labelled
// "c<charIndex>-l<lineIndex>" so a mismatch is easy to trace.
func buildValidBannerFixture() (lines []string, want Banner) {
	want = make(Banner, Count)
	for i := 0; i < Count; i++ {
		r := First + rune(i)
		lines = append(lines, "")
		block := make([]string, Height)
		for l := 0; l < Height; l++ {
			block[l] = fmt.Sprintf("c%d-l%d", i, l)
			lines = append(lines, block[l])
		}
		want[r] = block
	}
	return lines, want
}

// writeLines writes lines to name inside a fresh temp dir and returns the
// file's path.
func writeLines(t *testing.T, name string, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing fixture %s: %v", name, err)
	}
	return path
}

func TestLoad(t *testing.T) {
	validLines, validWant := buildValidBannerFixture()
	validPath := writeLines(t, "valid.txt", validLines)

	// One line short of Count*stride: the overall length check must
	// reject it before any block is even parsed.
	shortPath := writeLines(t, "short.txt", validLines[:len(validLines)-1])

	// Right number of lines, but the very first separator isn't blank.
	badSepLines := append([]string(nil), validLines...)
	badSepLines[0] = "not blank"
	badSepPath := writeLines(t, "badsep.txt", badSepLines)

	missingPath := filepath.Join(t.TempDir(), "missing.txt")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    Banner
		wantErr bool
	}{
		{"file does not exist", args{missingPath}, nil, true},
		{"wrong total line count", args{shortPath}, nil, true},
		{"missing blank separator", args{badSepPath}, nil, true},
		{"well-formed banner file", args{validPath}, validWant, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		sl []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty input", args{nil}, true},
		{"printable ascii", args{[]rune("Hello, World! 123")}, true},
		{"newline is allowed", args{[]rune("Hi\nThere")}, true},
		{"boundary runes space and tilde", args{[]rune{' ', '~'}}, true},
		{"control char below space", args{[]rune{31}}, false},
		{"tab is not allowed", args{[]rune{'\t'}}, false},
		{"DEL above tilde", args{[]rune{127}}, false},
		{"non-ascii rune", args{[]rune("héllo")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.args.sl); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
