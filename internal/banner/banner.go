// Package banner отвечает за загрузку и разбор файлов-баннеров (standard.txt,
// shadow.txt, thinkertoy.txt) и предоставляет доступ к ASCII-блоку (8 строк)
// для каждого печатного символа.
//
// Package banner is responsible for loading and parsing banner files
// (standard.txt, shadow.txt, thinkertoy.txt) and provides access to the
// 8-line ASCII block for each printable character.
package banner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// Height — сколько строк занимает один символ в баннере.
	// Height is how many lines a single character occupies.
	Height = 8

	// First и Last — границы диапазона символов, описанных в файле:
	// от пробела (32) до тильды (126) включительно.
	// First and Last are the bounds of the character range stored in a
	// banner file: from space (32) up to tilde (126), inclusive.
	First = ' '
	Last  = '~'

	// stride — размер одного блока в файле: пустая строка-разделитель + 8 строк.
	// stride is the size of one block in the file: a blank separator + 8 lines.
	stride = Height + 1

	// Count — количество описанных символов (95).
	// Count is the number of described characters (95).
	Count = int(Last-First) + 1
)

// Banner хранит для каждой печатной руны её ASCII-блок из Height строк.
// Banner maps every printable rune to its Height-line ASCII block.
type Banner map[rune][]string

// Load читает файл-баннер и раскладывает его на блоки по символам.
//
// Формат файла проверен по реальным standard.txt / shadow.txt / thinkertoy.txt:
// файл начинается с пустой строки, затем идут ровно 8 строк символа, затем
// снова пустая строка и т.д. Итого Count*stride == 855 строк.
//
// Load reads a banner file and splits it into per-character blocks. The file
// starts with a blank line, followed by exactly 8 lines for a character, then
// another blank line, and so on — Count*stride == 855 lines in total.
func Load(path string) (Banner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("banner: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Файлы могут прийти с CRLF — \r ломает выравнивание блоков.
		// Files may arrive with CRLF line endings; a stray \r breaks alignment.
		lines = append(lines, strings.TrimSuffix(scanner.Text(), "\r"))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("banner: reading %s: %w", path, err)
	}

	want := Count * stride
	if len(lines) != want {
		return nil, fmt.Errorf("banner: %s is malformed: got %d lines, want %d", path, len(lines), want)
	}

	b := make(Banner, Count)
	for i := 0; i < Count; i++ {
		// +1 — пропускаем пустую строку-разделитель перед блоком.
		// +1 skips the blank separator line that precedes each block.
		start := i*stride + 1
		if lines[start-1] != "" {
			return nil, fmt.Errorf("banner: %s is malformed: line %d must be a blank separator", path, start)
		}
		// Копируем срез: иначе блоки остались бы окнами в один общий массив,
		// и любая правка одного блока задела бы соседние.
		// Copy the slice: otherwise blocks would be windows into one shared
		// array and mutating one block would corrupt its neighbours.
		block := make([]string, Height)
		copy(block, lines[start:start+Height])
		b[First+rune(i)] = block
	}
	return b, nil
}

// Validate сообщает, состоит ли строка только из допустимых рун: печатный
// ASCII 32..126 либо перевод строки \n (10).
//
// Validate reports whether every rune is allowed: printable ASCII 32..126, or
// a newline \n (10).
func Validate(sl []rune) bool {
	for _, v := range sl {
		if v == '\n' {
			continue
		}
		if v < First || v > Last {
			return false
		}
	}
	return true
}
