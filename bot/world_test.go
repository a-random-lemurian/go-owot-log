package bot_test

import (
	"testing"

	"github.com/a-random-lemurian/go-owot-log/bot"
)

type testcase struct {
	Input string
	Hide bool
	Expected string
}

func TestOWOTUrlCreation(t *testing.T) {
	testcases := []testcase{
		{"", true, "wss://ourworldoftext.com/ws/?hide=1"},
		{"go-owot", false, "wss://ourworldoftext.com/go-owot/ws/"},
	}

	for _, testcase := range testcases {
		actual := bot.CreateOWOTUrl(testcase.Input, testcase.Hide)
		if actual != testcase.Expected {
			t.Errorf("did not get expected output: input %s, actual %s, expected %s", testcase.Input, actual, testcase.Expected)
		}
	}
}
