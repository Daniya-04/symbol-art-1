// Package printer отвечает за построение и вывод ASCII-арта из входной
// строки, используя баннер (banner.Banner) для символьных блоков.
//
// Package printer is responsible for building and printing ASCII art from
// an input string, using a banner.Banner for the character blocks.
package printer

// TODO (RU):
// 1. Реализовать основную функцию рендера, например:
//    func Render(input string, b banner.Banner) []string
//    которая возвращает готовые строки ASCII-арта (без вывода на экран),
//    чтобы main.go/analyzer могли и напечатать их, и посчитать размеры.
// 2. Разбить входную строку на строки по разделителю \n (strings.Split).
//    Продумать разницу между strings.Split и обработкой edge-кейсов:
//    - "" (пустая строка) -> вообще ничего не выводить (0 строк результата).
//    - "\n" -> ровно одна пустая строка вывода (8 строк пробелов? или
//      просто один перевод строки — свериться с примером в задании:
//      "\n" выводит одну пустую строку, а не 8).
//    - "Hello\n" -> арт для Hello (8 строк) + одна пустая строка.
//    - "Hello\n\nThere" -> арт Hello (8 строк) + одна пустая строка + арт There (8 строк).
//    Внимательно изучить пример из задания про "common mistake": нельзя
//    просто печатать одну пустую строку на каждый разделитель \n — нужно
//    посчитать точное количество разделителей и понять как оно соотносится
//    с количеством пустых строк.
// 3. Для каждой непустой строки (сегмента после Split) построить 8 строк
//    ASCII-арта:
//    - Для каждого символа в сегменте получить его 8-строчный блок через
//      banner.Get(r).
//    - Соединить блоки всех символов строки горизонтально: i-я строка
//      результата — это конкатенация i-х строк каждого блока символа
//      (i от 0 до 7).
// 4. Для строки-разделителя (когда встречается \n подряд без текста)
//    добавить одну пустую строку в результат (не 8 строк!).
// 5. Реализовать функцию вывода в stdout:
//    func Print(lines []string)
//    - Просто печатает каждую строку с переводом строки (fmt.Println или
//      bufio.Writer для эффективности при больших данных).
// 6. Реализовать вычисление размеров арта для блока "Art Dimensions":
//    func Dimensions(lines []string) (height int, width int)
//    - height — количество строк.
//    - width — длина самой длинной строки (или единой длины, если все
//      строки одной длины — уточнить по примерам, где ширина считается
//      по максимальной длине среди всех строк вывода).
//
// TODO (EN):
// 1. Implement the main render function, e.g.:
//    func Render(input string, b banner.Banner) []string
//    which returns the finished ASCII art lines (without printing them),
//    so that main.go/analyzer can both print them and compute dimensions.
// 2. Split the input string on \n (strings.Split). Think through the edge
//    cases carefully:
//    - "" (empty string) -> produce no output at all (0 result lines).
//    - "\n" -> exactly one empty output line (NOT 8 lines of spaces —
//      check the example in the spec: "\n" prints one blank line, not 8).
//    - "Hello\n" -> art for Hello (8 lines) + one blank line.
//    - "Hello\n\nThere" -> Hello art (8 lines) + one blank line + There art (8 lines).
//    Study the "common mistake" note in the spec closely: you must NOT
//    just print one blank line per \n separator — work out how the number
//    of separators actually maps to the number of blank lines produced.
// 3. For each non-empty segment (after Split), build 8 lines of ASCII art:
//    - For every character in the segment, fetch its 8-line block via
//      banner.Get(r).
//    - Concatenate all character blocks in the segment horizontally: the
//      i-th output line is the concatenation of the i-th line of every
//      character's block (i from 0 to 7).
// 4. For a separator case (consecutive \n with no text between them), add
//    exactly one blank line to the result (NOT 8 lines!).
// 5. Implement a function to print to stdout:
//    func Print(lines []string)
//    - Simply prints each line followed by a newline (fmt.Println or a
//      bufio.Writer for efficiency on larger inputs).
// 6. Implement art dimension calculation for the "Art Dimensions" block:
//    func Dimensions(lines []string) (height int, width int)
//    - height is the number of lines.
//    - width is the length of the longest line (or a single uniform
//      length if all lines share it — confirm against the examples, where
//      width is measured as the max length among all output lines).
