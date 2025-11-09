package bot_test

import (
	"testing"

	"github.com/a-random-lemurian/go-owot-log/bot"
)

type testcase struct {
	Input string
	Expected string
}

func TestOWOTUrlCreation(t *testing.T) {
	testcases := []testcase{
		{"", "wss://ourworldoftext.com/ws/"},
		{"go-owot", "wss://ourworldoftext.com/go-owot/ws/"},
	}

	for _, testcase := range testcases {
		actual := bot.CreateOWOTUrl(testcase.Input)
		if actual != testcase.Expected {
			t.Errorf("did not get expected output: input %s, actual %s, expected %s", testcase.Input, actual, testcase.Expected)
		}
	}
}
