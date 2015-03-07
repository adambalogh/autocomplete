package autocomp

import (
	"strings"
	"testing"
)

var newtests = []struct {
	Text string
	a    Autocomp
}{
	{
		"short text for short test for you",
		Autocomp{
			[]string{"short", "text", "for", "short", "test", "for", "you"},
			map[string]int{
				"short": 2,
				"text":  1,
				"for":   2,
				"test":  1,
				"you":   1,
			},
			map[string]Counter{
				"short": map[string]int{"text": 1, "test": 1},
				"text":  map[string]int{"for": 1},
				"for":   map[string]int{"short": 1, "you": 1},
				"test":  map[string]int{"for": 1},
				"you":   map[string]int{},
			},
		},
	},
}

func TestNew(t *testing.T) {
	for _, test := range newtests {
		a := New(strings.NewReader(test.Text))

		for word, count := range test.a.WordsCount {
			c, ok := a.WordsCount[word]
			if !ok {
				t.Errorf("Missing word: %s", word)
				continue
			}
			if c != count {
				t.Errorf("Expected count for '%s' is %d, got %d", word, count, c)
			}
		}

		for word, next := range test.a.WordTuples {
			for nextWord, count := range next {
				c, ok := a.WordTuples[word][nextWord]
				if !ok {
					t.Errorf("Not found: '%s' is followed by '%s'", word, nextWord)
					continue
				}
				if c != count {
					t.Errorf("Expected count for '%s' followed by '%s' is %d, got %d", word, nextWord, count, c)
				}
			}
		}
	}
}
