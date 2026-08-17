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
	"os"
)

type Banner map[rune][]string // Обьявили свой тип с навзанием Банер Мап - ключ - значение

func Load(path string) (Banner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	b := make(Banner)
	for i:=32; i< 127; i++{
		start := (i - 32) * 9
		end := start + 8
		b[rune(i)] = allLines[start:end]
	}

	return b, nil
	// TODO: здесь allLines уже содержит ВСЕ строки файла.
	// Дальше — ваша часть из пункта 3 TODO-списка ниже: пройтись по
	// allLines блоками по 9 строк (1 разделитель + 8 строк символа) и
	// заполнить Banner.
}

// func BannerCheck(input string) map[string][][]string {

// }

func Validate(sl []rune) bool {

	for _, v := range sl {
		if (v < 32 || v > 126) || v != 10 {
			return false
		}
	}
	return true
}

// 1. Определить тип для хранения баннера, например:
//    type Banner map[rune][]string
//    где значение — срез из 8 строк, представляющих символ.
// 2. Реализовать функцию загрузки файла баннера с диска, используя только
//    стандартный пакет os / io (без внешних библиотек), например:
//    func Load(path string) (Banner, error)
//    - Открыть файл, прочитать построчно.
//    - Баннер описывает символы по порядку начиная с пробела (код 32) и
//      далее по возрастанию ASCII-кода до кода 126.
//    - Каждый символ занимает ровно 8 строк, блоки разделены одной пустой
//      строкой (нужно свериться с реальным форматом входных файлов).
// 3. Продумать, как разбить файл на блоки по 8 строк и сопоставить каждый
//    блок с соответствующим руной (rune), начиная с ' ' (32) и увеличивая
//    код на 1 для каждого следующего блока.
// 4. Реализовать функцию получения блока символа:
//    func (b Banner) Get(r rune) []string
//    - Вернуть 8 строк для символа r.
//    - Продумать поведение для символов вне диапазона 32-126 (не должно
//      происходить, если Validate вызывается заранее, но стоит решить,
//      что возвращать на всякий случай).
// 5. Исправить/реализовать функцию Validate(sl []rune) bool:
//    - Должна возвращать true, если все руны — печатные ASCII (32-126)
//      ЛИБО символ новой строки \n.
//    - Текущая сигнатура использует v <= 32 || v >= 126, что неверно
//      (нужно понять, почему это неверно, и передумать условие: пробел
//      это код 32 и он ДОЛЖЕН быть валиден, а 126 '~' тоже валиден).
//    - Учесть, что \n (код 10) тоже должен считаться допустимым, так как
//      по условию входная строка может содержать \n.
// 6. Решить, где хранить путь к файлам баннеров по умолчанию (standard.txt
//    как банер по умолчанию) и как их встраивать в сборку (обычные файлы
//    рядом с go.mod, читаемые в рантайме).
//
// TODO (EN):
// 1. Define a type to store a banner, e.g.:
//    type Banner map[rune][]string
//    where the value is a slice of 8 strings representing one character.
// 2. Implement a function to load a banner file from disk using only the
//    standard os / io packages (no external libraries), e.g.:
//    func Load(path string) (Banner, error)
//    - Open the file, read it line by line.
//    - The banner describes characters in order starting from space
//      (code 32) and increasing by ASCII code up to code 126.
//    - Each character occupies exactly 8 lines, blocks are separated by a
//      blank line (verify this against the actual provided banner files).
// 3. Work out how to split the file into 8-line blocks and map each block
//    to its corresponding rune, starting at ' ' (32) and incrementing the
//    code by 1 for every following block.
// 4. Implement a function to fetch a character's block:
//    func (b Banner) Get(r rune) []string
//    - Return the 8 lines for character r.
//    - Decide what should happen for characters outside the 32-126 range
//      (should not happen if Validate is called first, but decide on a
//      safe fallback anyway).
// 5. Fix/implement the Validate(sl []rune) bool function:
//    - Should return true only if every rune is either printable ASCII
//      (32-126) OR a newline character \n.
//    - The current signature uses v <= 32 || v >= 126, which is incorrect
//      (work out why: space is code 32 and MUST be valid, and 126 '~' is
//      also valid, so re-think the condition).
//    - Remember that \n (code 10) must also be treated as valid, since the
//      input string is allowed to contain \n per the spec.
// 6. Decide where the default banner files live (standard.txt as the
//    default banner) and how they get read at runtime (plain files next to
//    go.mod, read at program startup).
