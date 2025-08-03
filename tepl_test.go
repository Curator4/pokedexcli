package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	} {
		{
			input: "   hello   world   ",
			expected: []string{"hello", "world"},
		},
		{
			input: "HELLO    WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input: "hello---world",
			expected: []string{"hello---world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		expected := c.expected

		if len(actual) != len(expected) {
			t.Errorf("len %v, expected %v", len(actual), len(expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := expected[i]

			if word != expectedWord {
				t.Errorf("got word %v, expected %v", word, expectedWord)
			}
		}
	}
}
