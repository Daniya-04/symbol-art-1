// Package printer отвечает за построение и вывод ASCII-арта из входной
// строки, используя баннер (banner.Banner) для символьных блоков.
//
// Package printer is responsible for building and printing ASCII art from
// an input string, using a banner.Banner for the character blocks.
package printer

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"symbol-art/internal/banner"
	"testing"
)

// testBlock builds a banner.Height-line block whose lines are labelled
// "<prefix>0".."<prefix>7", so the block a line came from (and its index)
// is visible directly in the rendered output.
func testBlock(prefix string) []string {
	lines := make([]string, banner.Height)
	for i := range lines {
		lines[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return lines
}

// testBanner is a small fake banner used across Render/renderSegment
// tests — real font files are far too big to hand-check expected output
// against.
func testBanner() banner.Banner {
	return banner.Banner{
		'A': testBlock("a"),
		'B': testBlock("b"),
	}
}

func TestRender(t *testing.T) {
	b := testBanner()
	type args struct {
		input string
		b     banner.Banner
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty input", args{"", b}, nil},
		{"single segment, no newline", args{"AB", b}, []string{"a0b0", "a1b1", "a2b2", "a3b3", "a4b4", "a5b5", "a6b6", "a7b7"}},
		{"two segments separated by one newline", args{"A\nB", b}, []string{
			"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7",
			"b0", "b1", "b2", "b3", "b4", "b5", "b6", "b7",
		}},
		{"blank line between segments", args{"A\n\nB", b}, []string{
			"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7",
			"",
			"b0", "b1", "b2", "b3", "b4", "b5", "b6", "b7",
		}},
		{"trailing newline adds one blank line", args{"A\n", b}, []string{
			"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Render(tt.args.input, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderSegment(t *testing.T) {
	b := testBanner()
	type args struct {
		segment string
		b       banner.Banner
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"single character", args{"A", b}, []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7"}},
		{"two characters joined horizontally", args{"AB", b}, []string{"a0b0", "a1b1", "a2b2", "a3b3", "a4b4", "a5b5", "a6b6", "a7b7"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderSegment(tt.args.segment, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

// captureStdout redirects os.Stdout to a pipe for the duration of fn and
// returns everything written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading captured stdout: %v", err)
	}
	return string(out)
}

func TestPrint(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no lines", args{nil}, ""},
		{"two lines", args{[]string{"foo", "bar"}}, "foo\nbar\n"},
		{"blank line in the middle", args{[]string{"", "x"}}, "\nx\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() { Print(tt.args.lines) })
			if got != tt.want {
				t.Errorf("Print() wrote %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDimensions(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name       string
		args       args
		wantHeight int
		wantWidth  int
	}{
		{"no lines", args{nil}, 0, 0},
		{"single empty line", args{[]string{""}}, 1, 0},
		{"varying widths", args{[]string{"Hi", "Hello there", ""}}, 3, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeight, gotWidth := Dimensions(tt.args.lines)
			if gotHeight != tt.wantHeight {
				t.Errorf("Dimensions() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("Dimensions() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
		})
	}
}
