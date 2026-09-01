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
	for _, r := range input {
		switch {
		case r == ' ' || r == '\n':
			b.Spaces++
		case r >= 'A' && r <= 'Z':
			b.Uppercase++
		case r >= 'a' && r <= 'z':
			b.Lowercase++
		case r >= '0' && r <= '9':
			b.Digits++
		default:
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
	p.MixedCase = hasUpper && hasLower

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
	seen := make(map[rune]bool)
	total := 0
	for _, r := range input {
		seen[r] = true
		total++
	}
	score := float64(len(seen)) / float64(total) * 100
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
	found := false
	if p.MixedCase {
		sb.WriteString("- Mixed case detected\n")
		found = true
	}
	for _, r := range p.RepeatedChars {
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

	hasUpper, hasLower, hasPunct := false, false, false
	for _, r := range input {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case strings.ContainsRune("!?.,;:", r):
			hasPunct = true
		}
	}

	if hasLower && !hasUpper {
		s = append(s, "Try an all-uppercase version for more visual impact: "+strings.ToUpper(input))
	}
	if hasUpper && !hasLower {
		s = append(s, "Try a lowercase version for a softer look: "+strings.ToLower(input))
	}
	if !hasPunct {
		s = append(s, `Consider adding punctuation (e.g. "!") to give it more character`)
	}

	switch words := strings.Fields(input); len(words) {
	case 0:
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
