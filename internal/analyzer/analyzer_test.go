// Package analyzer реализует rule-based анализ входного текста: подсчёт
// типов символов, обнаружение паттернов, вычисление показателя сложности
// (--analyze), а также генерацию текстовых предложений (--suggest).
//
// Package analyzer implements rule-based text analysis of the input
// string: character type counting, pattern detection, a complexity score
// (--analyze), and generating text suggestions (--suggest).
package analyzer

import (
	"reflect"
	"testing"
)

func TestClassifyChars(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want Breakdown
	}{
		{"empty input", args{""}, Breakdown{}},
		{"one of each category", args{"Ab3 !\n"}, Breakdown{Uppercase: 1, Lowercase: 1, Digits: 1, Special: 1, Spaces: 2}},
		{"digits only", args{"12345"}, Breakdown{Digits: 5}},
		{"special only", args{"!@#$%"}, Breakdown{Special: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyChars(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClassifyChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectPatterns(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want Patterns
	}{
		{"empty input", args{""}, Patterns{}},
		{"single char, nothing to detect", args{"a"}, Patterns{}},
		{"mixed case with a repeated run", args{"Hello"}, Patterns{MixedCase: true, RepeatedChars: []string{"ll"}}},
		{"lowercase only, multiple repeats and a digit run", args{"aabb1122"}, Patterns{RepeatedChars: []string{"aa", "bb", "11", "22"}, NumericSequences: []string{"1122"}}},
		{"mixed case with a numeric sequence, no repeats", args{"abc123XYZ"}, Patterns{MixedCase: true, NumericSequences: []string{"123"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectPatterns(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplexityScore(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"empty input", args{""}, 0},
		{"half unique", args{"aabbcc"}, 50},
		{"all unique", args{"AbCdEf"}, 100},
		{"case-sensitive uniqueness", args{"Aa"}, 100},
		{"single repeated char", args{"aaaa"}, 25},
		{"rounds to two decimals", args{"aaabbbccd"}, 44.44},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComplexityScore(tt.args.input); got != tt.want {
				t.Errorf("ComplexityScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatAnalysis(t *testing.T) {
	type args struct {
		input     string
		artHeight int
		artWidth  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty input, no patterns",
			args{"", 0, 0},
			"--- AI Analysis ---\n" +
				"Character Breakdown:\n" +
				"- Uppercase: 0\n" +
				"- Lowercase: 0\n" +
				"- Digits: 0\n" +
				"- Special characters: 0\n" +
				"- Spaces: 0\n" +
				"Patterns Detected:\n" +
				"- None detected\n" +
				"Complexity Score: 0.00%\n" +
				"Art Dimensions: 0 lines × 0 characters",
		},
		{
			"mixed case with a repeated run",
			args{"Hello!", 8, 30},
			"--- AI Analysis ---\n" +
				"Character Breakdown:\n" +
				"- Uppercase: 1\n" +
				"- Lowercase: 4\n" +
				"- Digits: 0\n" +
				"- Special characters: 1\n" +
				"- Spaces: 0\n" +
				"Patterns Detected:\n" +
				"- Mixed case detected\n" +
				"- Repeated characters: \"ll\"\n" +
				"Complexity Score: 83.33%\n" +
				"Art Dimensions: 8 lines × 30 characters",
		},
		{
			"numeric sequence, no mixed case or repeats",
			args{"ab12", 5, 10},
			"--- AI Analysis ---\n" +
				"Character Breakdown:\n" +
				"- Uppercase: 0\n" +
				"- Lowercase: 2\n" +
				"- Digits: 2\n" +
				"- Special characters: 0\n" +
				"- Spaces: 0\n" +
				"Patterns Detected:\n" +
				"- Numeric sequence: \"12\"\n" +
				"Complexity Score: 100.00%\n" +
				"Art Dimensions: 5 lines × 10 characters",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatAnalysis(tt.args.input, tt.args.artHeight, tt.args.artWidth); got != tt.want {
				t.Errorf("FormatAnalysis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuggest(t *testing.T) {
	type args struct {
		input     string
		artHeight int
		artWidth  int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"all lowercase, no punctuation, two words",
			args{"hello world", 8, 50},
			[]string{
				"Try an all-uppercase version for more visual impact: HELLO WORLD",
				`Consider adding punctuation (e.g. "!") to give it more character`,
				"This phrase has 2 words; a shorter word may render more compactly",
				"Output dimensions: 8 lines × 50 characters.",
			},
		},
		{
			"all uppercase, no punctuation, single word",
			args{"HELLO", 8, 25},
			[]string{
				"Try a lowercase version for a softer look: hello",
				`Consider adding punctuation (e.g. "!") to give it more character`,
				"This is a single word; try a short phrase for a richer banner",
				"Output dimensions: 8 lines × 25 characters.",
			},
		},
		{
			"mixed case with punctuation, double and edge spaces",
			args{"  Hi  There!  ", 3, 40},
			[]string{
				"This phrase has 2 words; a shorter word may render more compactly",
				"Multiple consecutive spaces detected; consider collapsing them",
				"Leading or trailing spaces detected; consider trimming them",
				"Output dimensions: 3 lines × 40 characters.",
			},
		},
		{
			"empty input",
			args{"", 0, 0},
			[]string{
				`Consider adding punctuation (e.g. "!") to give it more character`,
				"Output dimensions: 0 lines × 0 characters.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Suggest(tt.args.input, tt.args.artHeight, tt.args.artWidth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Suggest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSuggestions(t *testing.T) {
	type args struct {
		suggestions []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no suggestions", args{nil}, "--- AI Suggestions ---"},
		{"empty slice", args{[]string{}}, "--- AI Suggestions ---"},
		{"two suggestions", args{[]string{"a", "b"}}, "--- AI Suggestions ---\n- a\n- b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSuggestions(tt.args.suggestions); got != tt.want {
				t.Errorf("FormatSuggestions() = %v, want %v", got, tt.want)
			}
		})
	}
}
