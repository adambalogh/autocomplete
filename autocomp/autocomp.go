package autocomp

import (
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

// Counter is used to count the number of words.
type Counter map[string]int

// WordPrediction contains the predicted word, and it's count
// a higher count value means that the word is more likely.
type WordPrediction struct {
	Word  string
	Count int
}

// ByCount sorts WordPredictions by their Count value.
type ByCount []WordPrediction

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Count > c[j].Count }

// Autocomp is a Markov Chain model that can be used for autocompletion.
type Autocomp struct {
	Words      []string
	WordsCount Counter
	WordTuples map[string]Counter
}

// New returns a model that was trained on the text read from r.
func New(r io.Reader) *Autocomp {
	a := new(Autocomp)
	a.Words = make([]string, 0)
	a.WordsCount = make(map[string]int)
	a.WordTuples = make(map[string]Counter)

	w := regexp.MustCompile(`'?([a-zA-z'-]+)'?`)
	all, _ := ioutil.ReadAll(r)
	words := w.FindAllString(string(all), -1)

	for _, word := range words {
		word = strings.ToLower(word)
		a.Words = append(a.Words, word)
		a.WordsCount[word]++
	}

	for i := 0; i < len(a.Words)-1; i++ {
		first := a.Words[i]
		second := a.Words[i+1]
		if a.WordTuples[first] == nil {
			a.WordTuples[first] = make(map[string]int)
		}
		a.WordTuples[first][second]++
	}

	return a
}

// predictWord returns the top count words fron the wordCount map,
// that start with the string prefix.
func predictWord(wordsCount Counter, prefix string, count int) []WordPrediction {
	predictions := make([]WordPrediction, 0)
	for word, count := range wordsCount {
		if strings.HasPrefix(word, prefix) {
			p := WordPrediction{word, count}
			predictions = append(predictions, p)
		}
	}
	sort.Sort(ByCount(predictions))
	if len(predictions) > count {
		return predictions[:count]
	}
	return predictions
}

// Predict returns the most likely count words
//
// line could contain 1 or more words, the prediction will always
// be aimed at the last word, and if there are more than 1 words,
// the algorithm will use the previous word to predict better words.
func (a *Autocomp) Predict(line string, count int) []WordPrediction {
	words := strings.Fields(line)
	if len(words) < 2 {
		return predictWord(a.WordsCount, words[0], count)
	}
	return a.PredictNextWord(words[len(words)-2], words[len(words)-1], count)
}

// PredictNextWord returns the most likely word that has the prefix start,
// given that the previous word was first
func (a *Autocomp) PredictNextWord(first, second string, count int) []WordPrediction {
	return predictWord(a.WordTuples[first], second, count)
}
