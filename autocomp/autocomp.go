package autocomp

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Counter map[string]int

type WordPrediction struct {
	Word  string
	Count int
}

type ByCount []WordPrediction

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Count > c[j].Count }

type Autocomp struct {
	Words      []string
	WordsCount Counter
	WordTuples map[string]Counter
}

func New() *Autocomp {
	a := new(Autocomp)
	a.Words = make([]string, 0)
	a.WordsCount = make(map[string]int)
	a.WordTuples = make(map[string]Counter)

	big, err := os.Open("big.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	w := regexp.MustCompile(`'?([a-zA-z'-]+)'?`)
	all, _ := ioutil.ReadAll(big)
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

func (a *Autocomp) Predict(line string, count int) []WordPrediction {
	words := strings.Fields(line)
	if len(words) < 2 {
		return predictWord(a.WordsCount, words[0], count)
	}
	return a.PredictNextWord(words[len(words)-2], words[len(words)-1], count)
}

func (a *Autocomp) PredictNextWord(first, second string, count int) []WordPrediction {
	return predictWord(a.WordTuples[first], second, count)
}
