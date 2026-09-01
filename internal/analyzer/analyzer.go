// Package analyzer реализует rule-based анализ входного текста: подсчёт
// типов символов, обнаружение паттернов, вычисление показателя сложности
// (--analyze), а также генерацию текстовых предложений (--suggest).
//
// Package analyzer implements rule-based text analysis of the input
// string: character type counting, pattern detection, a complexity score
// (--analyze), and generating text suggestions (--suggest).
package analyzer

import (
	"fmt"
	"math"
	"strings"
)

// Breakdown хранит количество символов каждой категории во входной строке.
//
// Breakdown holds the count of each character category in the input string.
type Breakdown struct {
	Uppercase int
	Lowercase int
	Digits    int
	Special   int
	Spaces    int
}

// ClassifyChars классифицирует каждую руну input: заглавная/строчная буква,
// цифра, пробел (включая \n) или спецсимвол (всё остальное в 32-126).
//
// ClassifyChars classifies every rune of input: uppercase/lowercase letter,
// digit, whitespace (including \n), or special (everything else in 32-126).
func ClassifyChars(input string) Breakdown {
	var b Breakdown
	// range по string в Go идёт по РУНАМ (символам), а не по байтам —
	// это важно для не-ASCII текста, хотя banner.Validate у нас и так
	// не пускает дальше ничего, кроме печатного ASCII + '\n'.
	for _, r := range input {
		switch {
		// 'A'..'Z' и 'a'..'z' — это просто диапазоны кодов ASCII,
		// сравнение рун работает как сравнение чисел (rune — это
		// alias для int32), поэтому "r >= 'A' && r <= 'Z'" и есть
		// проверка "буква ли это, и заглавная ли".
		case r == ' ' || r == '\n':
			b.Spaces++
		case r >= 'A' && r <= 'Z':
			b.Uppercase++
		case r >= 'a' && r <= 'z':
			b.Lowercase++
		case r >= '0' && r <= '9':
			b.Digits++
		default:
			// Всё, что не попало ни в один case выше, но прошло
			// Validate (то есть печатный ASCII 32-126) — это
			// пунктуация и символы вроде !@#$% и т.п.
			b.Special++
		}
	}
	return b
}

// Patterns хранит результаты обнаружения паттернов во входной строке.
//
// Patterns holds the pattern-detection results for the input string.
type Patterns struct {
	MixedCase        bool
	RepeatedChars    []string
	NumericSequences []string
}

// DetectPatterns ищет смешанный регистр, повторяющиеся подряд символы (2+)
// и числовые последовательности (2+ цифры подряд).
//
// DetectPatterns looks for mixed case, runs of 2+ identical characters in a
// row, and numeric sequences of 2+ digits in a row.
func DetectPatterns(input string) Patterns {
	// Переводим строку в []rune один раз, чтобы дальше можно было
	// индексировать по позиции символа (runes[i]) — со string так
	// напрямую нельзя, если в ней могут быть многобайтовые символы.
	runes := []rune(input)

	var p Patterns
	hasUpper, hasLower := false, false
	for _, r := range runes {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		}
		if r >= 'a' && r <= 'z' {
			hasLower = true
		}
	}
	// Mixed case = встретилась хотя бы одна заглавная И хотя бы одна
	// строчная буква где угодно в строке (не обязательно подряд).
	p.MixedCase = hasUpper && hasLower

	// Поиск повторов подряд — классический "two pointer" проход:
	// i — начало текущего повторяющегося блока, j — бежит вперёд,
	// пока символ не изменится. Если блок [i:j) длиннее 1 символа —
	// значит нашли повтор ("oo", "ll" и т.п.), и i перескакивает
	// сразу на j (а не на i+1), чтобы не находить один и тот же
	// повтор много раз частично.
	for i := 0; i < len(runes); {
		j := i + 1
		for j < len(runes) && runes[j] == runes[i] {
			j++
		}
		if j-i >= 2 {
			p.RepeatedChars = append(p.RepeatedChars, string(runes[i:j]))
		}
		i = j
	}

	// Та же идея, но для цифр: ищем подряд идущие runes[i]..runes[j-1],
	// которые все являются цифрами '0'-'9'. Не-цифры просто
	// пропускаем по одному символу (i++).
	for i := 0; i < len(runes); {
		if runes[i] < '0' || runes[i] > '9' {
			i++
			continue
		}
		j := i + 1
		for j < len(runes) && runes[j] >= '0' && runes[j] <= '9' {
			j++
		}
		if j-i >= 2 {
			p.NumericSequences = append(p.NumericSequences, string(runes[i:j]))
		}
		i = j
	}

	return p
}

// ComplexityScore вычисляет (Unique_Characters / Total_Characters) * 100,
// подсчёт регистрозависимый, округление до 2 знаков после запятой.
//
// ComplexityScore computes (Unique_Characters / Total_Characters) * 100,
// case-sensitive, rounded to 2 decimal places.
func ComplexityScore(input string) float64 {
	if input == "" {
		return 0
	}
	// map[rune]bool как "множество" (set) — самый простой способ в Go
	// посчитать УНИКАЛЬНЫЕ символы: каждый r кладём ключом в map,
	// повторы просто перезатирают тот же ключ, поэтому len(seen) в
	// итоге — это и есть количество разных символов.
	//
	// Регистрозависимость получается "бесплатно": 'A' и 'a' — это
	// разные rune (разные числа), поэтому map их не схлопывает, и
	// "Aa" даёт 2 уникальных символа из 2 (100%), а не 1 из 2 (50%),
	// как было бы при регистронезависимом сравнении.
	seen := make(map[rune]bool)
	total := 0
	for _, r := range input {
		seen[r] = true
		total++
	}
	score := float64(len(seen)) / float64(total) * 100
	// Округление до 2 знаков через math.Round: умножаем на 100,
	// округляем до целого, делим обратно — стандартный приём, т.к.
	// у math нет отдельной функции "округли до N знаков".
	return math.Round(score*100) / 100
}

// FormatAnalysis форматирует блок "--- AI Analysis ---" для input,
// используя artHeight/artWidth для строки Art Dimensions.
//
// FormatAnalysis formats the "--- AI Analysis ---" block for input, using
// artHeight/artWidth for the Art Dimensions line.
func FormatAnalysis(input string, artHeight, artWidth int) string {
	b := ClassifyChars(input)
	p := DetectPatterns(input)
	score := ComplexityScore(input)

	var sb strings.Builder
	sb.WriteString("--- AI Analysis ---\n")
	sb.WriteString("Character Breakdown:\n")
	fmt.Fprintf(&sb, "- Uppercase: %d\n", b.Uppercase)
	fmt.Fprintf(&sb, "- Lowercase: %d\n", b.Lowercase)
	fmt.Fprintf(&sb, "- Digits: %d\n", b.Digits)
	fmt.Fprintf(&sb, "- Special characters: %d\n", b.Special)
	fmt.Fprintf(&sb, "- Spaces: %d\n", b.Spaces)

	sb.WriteString("Patterns Detected:\n")
	// found отслеживает, вывели ли мы хоть одну строку про паттерны —
	// если ни один из трёх видов паттернов не сработал, печатаем
	// "None detected" вместо пустого списка.
	found := false
	if p.MixedCase {
		sb.WriteString("- Mixed case detected\n")
		found = true
	}
	for _, r := range p.RepeatedChars {
		// %q оборачивает строку в кавычки и экранирует спецсимволы —
		// удобно для вывода вида Repeated characters: "ll".
		fmt.Fprintf(&sb, "- Repeated characters: %q\n", r)
		found = true
	}
	for _, n := range p.NumericSequences {
		fmt.Fprintf(&sb, "- Numeric sequence: %q\n", n)
		found = true
	}
	if !found {
		sb.WriteString("- None detected\n")
	}

	fmt.Fprintf(&sb, "Complexity Score: %.2f%%\n", score)
	fmt.Fprintf(&sb, "Art Dimensions: %d lines × %d characters", artHeight, artWidth)

	return sb.String()
}

// Suggest генерирует список текстовых предложений на основе регистра,
// пунктуации, количества слов и пробельных паттернов input. Последняя
// строка всегда — размеры вывода.
//
// Suggest generates a list of text suggestions based on input's case,
// punctuation, word count, and whitespace patterns. The last line is
// always the output dimensions.
func Suggest(input string, artHeight, artWidth int) []string {
	var s []string

	// Один проход по строке — сразу собираем три флага, которые
	// дальше решают, какие советы добавлять.
	hasUpper, hasLower, hasPunct := false, false, false
	for _, r := range input {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case strings.ContainsRune("!?.,;:", r):
			// strings.ContainsRune проверяет, есть ли r среди
			// перечисленных символов пунктуации — короче, чем
			// switch/case на каждый знак препинания отдельно.
			hasPunct = true
		}
	}

	// Регистр: если строка целиком строчная — предлагаем UPPER,
	// если целиком заглавная — предлагаем lower. Если уже смешанная
	// (hasUpper && hasLower), не предлагаем ничего про регистр.
	if hasLower && !hasUpper {
		s = append(s, "Try an all-uppercase version for more visual impact: "+strings.ToUpper(input))
	}
	if hasUpper && !hasLower {
		s = append(s, "Try a lowercase version for a softer look: "+strings.ToLower(input))
	}
	if !hasPunct {
		s = append(s, `Consider adding punctuation (e.g. "!") to give it more character`)
	}

	// strings.Fields бьёт строку по пробельным символам и
	// автоматически схлопывает несколько пробелов подряд, то есть
	// возвращает именно "слова", а не сырые куски между пробелами.
	switch words := strings.Fields(input); len(words) {
	case 0:
		// input состоит только из пробелов/пусто — про количество
		// слов советовать нечего.
	case 1:
		s = append(s, "This is a single word; try a short phrase for a richer banner")
	default:
		s = append(s, fmt.Sprintf("This phrase has %d words; a shorter word may render more compactly", len(words)))
	}

	if strings.Contains(input, "  ") {
		s = append(s, "Multiple consecutive spaces detected; consider collapsing them")
	}
	if strings.HasPrefix(input, " ") || strings.HasSuffix(input, " ") {
		s = append(s, "Leading or trailing spaces detected; consider trimming them")
	}

	// Эта строка добавляется ВСЕГДА и последней — по заданию блок
	// suggestions обязан заканчиваться размерами вывода.
	s = append(s, fmt.Sprintf("Output dimensions: %d lines × %d characters.", artHeight, artWidth))

	return s
}

// FormatSuggestions форматирует блок "--- AI Suggestions ---" со списком
// suggestions, каждое с префиксом "- ".
//
// FormatSuggestions formats the "--- AI Suggestions ---" block with the
// suggestions list, each prefixed with "- ".
func FormatSuggestions(suggestions []string) string {
	var sb strings.Builder
	sb.WriteString("--- AI Suggestions ---")
	for _, s := range suggestions {
		sb.WriteString("\n- ")
		sb.WriteString(s)
	}
	return sb.String()
}
