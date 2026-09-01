// Package printer отвечает за построение и вывод ASCII-арта из входной
// строки, используя баннер (banner.Banner) для символьных блоков.
//
// Package printer is responsible for building and printing ASCII art from
// an input string, using a banner.Banner for the character blocks.
package printer

import (
	"reflect"
	"symbol-art/internal/banner"
	"testing"
)

func TestRender(t *testing.T) {
	type args struct {
		input string
		b     banner.Banner
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
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
	type args struct {
		segment string
		b       banner.Banner
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderSegment(tt.args.segment, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(tt.args.lines)
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
		// TODO: Add test cases.
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
