// Package printer отвечает за построение и вывод ASCII-арта из входной
// строки, используя баннер (banner.Banner) для символьных блоков.
//
// Package printer is responsible for building and printing ASCII art from
// an input string, using a banner.Banner for the character blocks.
package printer

import (
	"fmt"
	"strings"

	"symbol-art/internal/banner"
)

// Render превращает input в готовые строки ASCII-арта, используя b для
// получения блока каждого символа. Пустая строка даёт nil (нет вывода).
// Каждый \n даёт одну пустую строку результата вместо Height строк.
//
// Render turns input into finished ASCII-art lines, using b to look up
// each character's block. An empty string yields nil (no output). Each \n
// yields a single blank result line instead of Height lines.
func Render(input string, b banner.Banner) []string {
	if input == "" {
		return nil
	}

	var lines []string
	for _, segment := range strings.Split(input, "\n") {
		if segment == "" {
			lines = append(lines, "")
			continue
		}
		lines = append(lines, renderSegment(segment, b)...)
	}
	return lines
}

// renderSegment строит Height строк для одного сегмента текста (без \n),
// склеивая по горизонтали блоки всех его символов.
//
// renderSegment builds Height lines for one text segment (no \n), joining
// the blocks of all its characters horizontally.
func renderSegment(segment string, b banner.Banner) []string {
	rows := make([]strings.Builder, banner.Height)
	for _, r := range segment {
		block := b[r]
		for i := 0; i < banner.Height; i++ {
			rows[i].WriteString(block[i])
		}
	}

	out := make([]string, banner.Height)
	for i := range rows {
		out[i] = rows[i].String()
	}
	return out
}

// Print печатает каждую строку в stdout.
//
// Print prints every line to stdout.
func Print(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

// Dimensions возвращает высоту (число строк) и ширину (длину самой
// длинной строки) готового арта.
//
// Dimensions returns the height (number of lines) and width (length of the
// longest line) of the rendered art.
func Dimensions(lines []string) (height, width int) {
	height = len(lines)
	for _, line := range lines {
		if n := len(line); n > width {
			width = n
		}
	}
	return height, width
}
