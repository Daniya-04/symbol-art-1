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
	// Пустая строка — особый случай: нет вообще никакого вывода
	// (ни одной строки, даже пустой).
	if input == "" {
		return nil
	}

	var lines []string
	// strings.Split режет строку по каждому вхождению "\n" и всегда
	// возвращает N+1 элементов на N разделителей. Например:
	//   "Hello"          -> ["Hello"]              (1 сегмент)
	//   "Hello\nThere"    -> ["Hello", "There"]      (2 сегмента, между
	//                        ними НЕТ пустой строки — просто два блока
	//                        друг под другом)
	//   "Hello\n\nThere"  -> ["Hello", "", "There"]  (средний сегмент
	//                        пустой — вот он и даёт ту самую одну
	//                        пустую строку между блоками)
	//   "Hello\n"         -> ["Hello", ""]           (пустой сегмент в
	//                        конце — блок Hello + одна пустая строка)
	for _, segment := range strings.Split(input, "\n") {
		if segment == "" {
			// Пустой сегмент — это НЕ повод рисовать 8-строчный блок
			// (получится 8 пустых строк вместо одной). Просто
			// добавляем одну пустую строку в результат.
			lines = append(lines, "")
			continue
		}
		// Непустой сегмент — обычный текст без переводов строки,
		// его рисуем как один блок высотой banner.Height.
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
	// Идея горизонтальной склейки: у нас Height (=8) "полос", и для
	// каждого символа сегмента мы дописываем его i-ю строку блока в
	// конец i-й полосы. В итоге rows[0] — это первая строка ВСЕХ
	// символов подряд, rows[1] — вторая строка всех символов, и т.д.
	// Например для "Hi": rows[0] = block['H'][0] + block['i'][0].
	//
	// strings.Builder используется вместо конкатенации через "+",
	// чтобы не пересоздавать строку на каждый символ — на длинных
	// строках это заметно быстрее.
	rows := make([]strings.Builder, banner.Height)
	for _, r := range segment {
		// b — это map[rune][]string, поэтому b[r] сразу даёт готовый
		// блок из Height строк для символа r (banner.Validate уже
		// гарантирует, что символ поддерживается, до вызова Render).
		block := b[r]
		for i := 0; i < banner.Height; i++ {
			rows[i].WriteString(block[i])
		}
	}

	// Превращаем Height builder'ов обратно в Height обычных строк.
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
	// Высота — это просто количество строк в готовом арте (включая
	// пустые строки-разделители, они тоже часть вывода).
	height = len(lines)
	// Ширина — длина самой длинной строки. Строки разной длины
	// возможны, например, если в input несколько сегментов из разных
	// по количеству символов слов ("Hi\nHello there").
	for _, line := range lines {
		if n := len(line); n > width {
			width = n
		}
	}
	return height, width
}
