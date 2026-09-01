# symbol-art

## RU

### 1. Цель проекта

Программа переводит текстовую строку в ASCII-арт, используя файл-баннер
(шрифт), а также умеет анализировать текст (`--analyze`) и предлагать
альтернативы форматирования (`--suggest`).

### 2. Правила обработки

- Формат баннер-файла: 8 строк на символ (константа `banner.Height`),
  символы описаны в диапазоне печатного ASCII `32..126` (пробел — тильда).
- Пустая строка `""` — нет вывода вообще (0 строк).
- `\n` — разделитель строк: непустой сегмент между разделителями даёт
  8 строк ASCII-арта, пустой сегмент — одну пустую строку вывода.
- Поддерживаемые символы: печатный ASCII `32-126` и `\n` (`banner.Validate`
  отклоняет всё остальное, например не-ASCII буквы или управляющие коды).

### 3. Использование AI-функций и примеры

`--analyze` — печатает блок с разбивкой символов по категориям
(заглавные/строчные/цифры/спецсимволы/пробелы), обнаруженные паттерны
(смешанный регистр, повторяющиеся символы, числовые последовательности),
показатель сложности и размеры арта:

```
$ go run . --analyze "Hello2024oo"
...ASCII-арт...

--- AI Analysis ---
Character Breakdown:
- Uppercase: 1
- Lowercase: 6
- Digits: 4
- Special characters: 0
- Spaces: 0
Patterns Detected:
- Mixed case detected
- Repeated characters: "ll"
- Repeated characters: "oo"
- Numeric sequence: "2024"
Complexity Score: 63.64%
Art Dimensions: 8 lines × 81 characters
```

`--suggest` — печатает блок с текстовыми предложениями (регистр,
пунктуация, количество слов, пробельные паттерны) и размерами вывода:

```
$ go run . --suggest "hello world"
...ASCII-арт...

--- AI Suggestions ---
- Try an all-uppercase version for more visual impact: HELLO WORLD
- Consider adding punctuation (e.g. "!") to give it more character
- This phrase has 2 words; a shorter word may render more compactly
- Output dimensions: 8 lines × 75 characters.
```

Это rule-based анализ: никаких внешних библиотек или ML — только простые
условия и подсчёты над рунами строки (пакет `internal/analyzer`).

### 4. Установка и запуск

Требования: Go версии из [go.mod](go.mod) (`go 1.26.4`), только стандартная
библиотека — внешних зависимостей нет.

```
git clone https://01.tomorrow-school.ai/git/abaitas/symbol-art.git
cd symbol-art
go run . [--analyze] [--suggest] [--banner=standard|shadow|thinkertoy] "STRING"
# либо собрать бинарник:
go build -o symbol-art .
./symbol-art "Hello"
```

Тесты:

```
go test ./...
```

### 5. Структура проекта

- [main.go](main.go) — точка входа: парсит флаги, загружает баннер,
  вызывает рендер и (опционально) блоки анализа/предложений.
- [internal/banner](internal/banner) — загрузка и валидация файлов-баннеров
  (`Load`, `Validate`).
- [internal/printer](internal/printer) — построение ASCII-арта из строки и
  баннера (`Render`), вывод (`Print`), вычисление размеров (`Dimensions`).
- [internal/analyzer](internal/analyzer) — rule-based анализ текста и
  генерация предложений (`--analyze`, `--suggest`).
- [standard.txt](standard.txt), [shadow.txt](shadow.txt),
  [thinkertoy.txt](thinkertoy.txt) — файлы-баннеры (шрифты).

## EN

### 1. Project objective

The program converts a text string into ASCII art using a banner file
(font), and can also analyze the text (`--analyze`) and suggest formatting
alternatives (`--suggest`).

### 2. Processing rules

- Banner file format: 8 lines per character (`banner.Height`), characters
  cover the printable ASCII range `32..126` (space through tilde).
- An empty string `""` — no output at all (0 lines).
- `\n` is a line separator: a non-empty segment between separators produces
  8 lines of ASCII art, an empty segment produces a single blank output
  line.
- Supported characters: printable ASCII `32-126` and `\n` (`banner.Validate`
  rejects everything else, e.g. non-ASCII letters or control codes).

### 3. AI feature usage and examples

`--analyze` prints a block with a character-category breakdown
(uppercase/lowercase/digits/special/spaces), detected patterns (mixed case,
repeated characters, numeric sequences), a complexity score, and the art
dimensions — see the example above.

`--suggest` prints a block of text suggestions (case, punctuation, word
count, whitespace patterns) plus the output dimensions — see the example
above.

This is rule-based analysis: no external libraries or ML, just plain
conditionals and counting over the string's runes (`internal/analyzer`
package).

### 4. Installation and setup

Requirements: the Go version from [go.mod](go.mod) (`go 1.26.4`), standard
library only — no external dependencies.

```
git clone https://01.tomorrow-school.ai/git/abaitas/symbol-art.git
cd symbol-art
go run . [--analyze] [--suggest] [--banner=standard|shadow|thinkertoy] "STRING"
# or build a binary:
go build -o symbol-art .
./symbol-art "Hello"
```

Running tests:

```
go test ./...
```

### 5. Folder structure

- [main.go](main.go) — entry point: parses flags, loads the banner, calls
  the renderer and (optionally) the analysis/suggestions blocks.
- [internal/banner](internal/banner) — loading and validating banner files
  (`Load`, `Validate`).
- [internal/printer](internal/printer) — building ASCII art from a string
  and a banner (`Render`), printing (`Print`), computing dimensions
  (`Dimensions`).
- [internal/analyzer](internal/analyzer) — rule-based text analysis and
  suggestion generation (`--analyze`, `--suggest`).
- [standard.txt](standard.txt), [shadow.txt](shadow.txt),
  [thinkertoy.txt](thinkertoy.txt) — banner files (fonts).
