# symbol-art-1

<!--
TODO (RU): README должен содержать следующие разделы (заполнить после
реализации проекта, ничего не выдумывать заранее — описывать то, что
реально реализовано):

1. Project objective / Цель проекта
   - Кратко: программа переводит текстовую строку в ASCII-арт, используя
     файл-баннер, и умеет анализировать текст (--analyze) и предлагать
     альтернативы (--suggest).

2. Processing rules / Правила обработки
   - Формат баннер-файла: 8 строк на символ.
   - Поведение при пустой строке "" — нет вывода вообще.
   - Поведение \n как разделителя строк (одна пустая строка на \n, не 8).
   - Диапазон поддерживаемых символов: печатные ASCII 32-126 + \n.

3. AI feature(s) usage and examples / Использование AI-функций и примеры
   - Пример команды с --analyze и пример вывода.
   - Пример команды с --suggest и пример вывода.
   - Пояснить, что это rule-based анализ (без внешних библиотек/ML).

4. Installation and setup / Установка и запуск
   - Требования: Go (указать версию из go.mod).
   - Команды: git clone, go build / go run .
   - Как запускать тесты: go test ./...

5. Folder structure / Структура проекта
   - Описать назначение main.go, internal/banner, internal/printer,
     internal/analyzer, файлов *.txt (баннеры), samples/.

TODO (EN): The README must contain the following sections (fill in after
the project is actually implemented — do not invent details ahead of
time, describe what is actually implemented):

1. Project objective
   - Briefly: the program converts a text string into ASCII art using a
     banner file, and can analyze the text (--analyze) and suggest
     alternatives (--suggest).

2. Processing rules
   - Banner file format: 8 lines per character.
   - Behavior for an empty string "" — no output at all.
   - Behavior of \n as a line separator (one blank line per \n, not 8).
   - Supported character range: printable ASCII 32-126 + \n.

3. AI feature(s) usage and examples
   - An example command using --analyze and its example output.
   - An example command using --suggest and its example output.
   - Clarify that this is rule-based analysis (no external libraries/ML).

4. Installation and setup
   - Requirements: Go (state the version from go.mod).
   - Commands: git clone, go build / go run .
   - How to run tests: go test ./...

5. Folder structure
   - Describe the purpose of main.go, internal/banner, internal/printer,
     internal/analyzer, the *.txt banner files, and samples/.
-->
