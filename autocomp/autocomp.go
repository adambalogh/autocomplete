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

// WordPrediction is a prediction for a word, where
// higher Count means higher probability.
type WordPrediction struct {
	Word  string
	Count int
}

// ByCount sorts WordPredictions by their probability.
type ByCount []WordPrediction

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Count > c[j].Count }

// Autocomp is a Markov Chain model that can be used for predicting
// what word is being typed.
type Autocomp struct {
	// Words contains all the words in the model
	Words []string
	// WordsCount counts each word
	WordsCount Counter
	// WordTuples counts how many times each word is followed by a specific word
	//
	// e.g. "this is" -> WordTuples["this"]["is"] == 1
	WordTuples map[string]Counter
}

// New returns a model that was trained on the text read from r.
func New(r io.Reader) *Autocomp {
	a := new(Autocomp)
	a.Words = make([]string, 0)
	a.WordsCount = make(Counter)
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
			a.WordTuples[first] = make(Counter)
		}
		a.WordTuples[first][second]++
	}

	return a
}

// predictWord returns the words with the highest probability,
// that start with the given prefix.
func predictWord(wordsCount Counter, prefix string, count int) []WordPrediction {
	predictions := make([]WordPrediction, 0)
	for word, count := range wordsCount {
		if strings.HasPrefix(word, prefix) {
			p := WordPrediction{word, count}
			predictions = append(predictions, p)
		}
	}
	sort.Sort(ByCount(predictions))
	if len(predictions) >= count {
		return predictions[:count]
	}
	return predictions
}

// Predict returns the most likely words the user is typing.
//
// line could contain 1 or more words, the prediction will always
// be aimed at the last word. If there are more than 1 words given,
// the algorithm will use the last but one word to predict the current word.
func (a *Autocomp) Predict(line string, count int) []WordPrediction {
	words := strings.Fields(line)
	if len(words) < 2 {
		return predictWord(a.WordsCount, words[0], count)
	}
	return a.PredictNextWord(words[len(words)-2], words[len(words)-1], count)
}

// PredictNextWord returns the most likely word that has the given prefix,
// given the previous word
func (a *Autocomp) PredictNextWord(first, second string, count int) []WordPrediction {
	return predictWord(a.WordTuples[first], second, count)
}
