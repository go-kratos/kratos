package gse

import (
	"fmt"
	"testing"
)

func printTokens(tokens []*Token, numTokens int) (output string) {
	for iToken := 0; iToken < numTokens; iToken++ {
		for _, word := range tokens[iToken].text {
			output += fmt.Sprint(string(word))
		}
		output += " "
	}
	return
}

func toWords(strings ...string) []Text {
	words := []Text{}
	for _, s := range strings {
		words = append(words, []byte(s))
	}
	return words
}

func bytesToString(bytes []Text) (output string) {
	for _, b := range bytes {
		output += (string(b) + "/")
	}
	return
}

func expect(t *testing.T, expect string, actual interface{}) {
	actualString := fmt.Sprint(actual)
	if expect != actualString {
		t.Errorf("期待值=\"%s\", 实际=\"%s\"", expect, actualString)
	}
}
