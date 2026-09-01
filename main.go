package main

import (
	"flag"
	"fmt"
	"os"

	"symbol-art/internal/analyzer"
	"symbol-art/internal/banner"
	"symbol-art/internal/printer"
)

func main() {
	// flag.Bool/flag.String регистрируют флаги и возвращают *bool/*string —
	// сами значения появятся только ПОСЛЕ flag.Parse(), поэтому дальше по
	// коду везде используется разыменование (*analyze, *bannerName).
	analyze := flag.Bool("analyze", false, "print an AI-style analysis of the input")
	suggest := flag.Bool("suggest", false, "print AI-style suggestions for the input")
	bannerName := flag.String("banner", "standard", "banner font to use: standard, shadow, or thinkertoy")
	flag.Parse()

	// flag.Parse() останавливается на первом аргументе без "-"/"--" —
	// всё, что осталось после него, это flag.Args(). У нас это должен
	// быть ровно один элемент: сама строка для рендера. Поэтому флаги
	// нужно указывать ДО строки: `go run . --analyze "Hello"`, а не
	// после.
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: go run . [--analyze] [--suggest] [--banner=standard|shadow|thinkertoy] STRING")
		os.Exit(1)
	}
	input := args[0]

	// Проверяем ДО загрузки баннера: нет смысла читать файл шрифта,
	// если строка всё равно содержит недопустимые символы.
	if !banner.Validate([]rune(input)) {
		fmt.Fprintln(os.Stderr, "Error: input contains unsupported characters")
		os.Exit(1)
	}

	// Имя шрифта превращаем в путь к файлу: "standard" -> "standard.txt".
	// Если файла с таким именем нет (опечатка в --banner или он не в
	// текущей директории), banner.Load вернёт понятную ошибку.
	b, err := banner.Load(*bannerName + ".txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Render строит сами строки арта, Print — просто выводит их в stdout.
	// Разделены специально: lines ещё нужны ниже для Dimensions.
	lines := printer.Render(input, b)
	printer.Print(lines)

	// Блоки --analyze и --suggest считаются от ОРИГИНАЛЬНОЙ строки
	// input, а не от отрендеренного арта — единственное, что из арта
	// им нужно, это его размеры (height/width) для последней строки
	// вывода.
	if *analyze || *suggest {
		height, width := printer.Dimensions(lines)
		if *analyze {
			fmt.Println()
			fmt.Println(analyzer.FormatAnalysis(input, height, width))
		}
		if *suggest {
			fmt.Println()
			fmt.Println(analyzer.FormatSuggestions(analyzer.Suggest(input, height, width)))
		}
	}
}
