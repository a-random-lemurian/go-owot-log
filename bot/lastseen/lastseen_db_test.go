package lastseen_test

import (
	"testing"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/lastseen"

	"github.com/a-random-lemurian/go-owot"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// scenarios
// 1. user gets a consenting user's lastseen
// 2. user gets a non-consenting user's lastseen
// 3. non-consenting user gets their own lastseen

func setup() (*lastseen.GormLastseenDB, error) {
	memdb, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	return lastseen.NewGormLastseenDB(memdb), nil
}

func makeMsg(realUsername string, t time.Time) *owot.MessageChat {
	return &owot.MessageChat{Kind: "chat", RealUsername: realUsername, Date: owot.UnixMillis(t)}
}

type lastseenTestcase struct {
	user  string
	milli int64
}

func TestUserGetsLastseen(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	lastseenTcases := []lastseenTestcase{
		{"privacyguy", 690000},
		{"otherentity", 1480000},
		{"funnierguy", 6942069},
	}
	for _, tcase := range lastseenTcases {
		err := db.SetLastseen(tcase.user, time.UnixMilli(tcase.milli))
		if err != nil {
			t.Errorf("error while setting up db: %e", err)
		}
	}

	for _, tcase := range lastseenTcases {
		t.Run(tcase.user, func(t *testing.T) {
			lastseen, err := db.GetLastseen(tcase.user)
			if err != nil {
				t.Error(err)
			}
			if lastseen.UnixMilli() != tcase.milli {
				t.Errorf("incorrect time for %s: got %d expected %d",
					tcase.user, lastseen.UnixMilli(), tcase.milli)
			}
		})
	}
}

type consentTestcase struct {
	user    string
	consent bool
	desc    string
}

func TestConsent(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = db.SetConsent("dontlogme", false)
	if err != nil {
		t.Error(err)
	}

	err = db.SetConsent("surelogme", true)
	if err != nil {
		t.Error(err)
	}

	consentTestcases := []consentTestcase{
		{"dontlogme", false, "opt out"},
		{"surelogme", true, "explicit consent"},
		{"idontexist", true, "not opted out"},
	}

	for _, tcase := range consentTestcases {
		t.Run(tcase.desc, func(t *testing.T) {
			consent := db.Consents(tcase.user)
			if consent != tcase.consent {
				t.Errorf("wrong consent for %s : should be %t", tcase.user, tcase.consent)
			}
		})
	}
}

func TestDefaultConsent(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = db.SetLastseen("privacyguy", time.UnixMilli(1000))
	if err != nil {
		t.Error(err)
	}

	if !db.Consents("privacyguy") {
		t.Errorf("consent needs to be 1 by default!!")
	}
}

func TestLaterOptout(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = db.SetLastseen("privacyguy", time.UnixMilli(1000))
	if err != nil {
		t.Error(err)
	}

	err = db.SetConsent("privacyguy", false)
	if err != nil {
		t.Error(err)
	}

	err = db.SetLastseen("privacyguy", time.UnixMilli(1005))

	if err != nil {
		t.Error(err)
	}

	if db.Consents("privacyguy") {
		t.Errorf("consent needs to be 0!!")
	}
}

func TestPreemptiveOptout(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = db.SetConsent("privacyguy", false)
	if err != nil {
		t.Error(err)
	}

	err = db.SetLastseen("privacyguy", time.UnixMilli(1005))

	if err != nil {
		t.Error(err)
	}

	if db.Consents("privacyguy") {
		t.Errorf("consent needs to be 0!!")
	}
}

func TestPreemptiveOptin(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = db.SetConsent("logme", true)
	if err != nil {
		t.Error(err)
	}

	err = db.SetLastseen("logme", time.UnixMilli(1005))

	if err != nil {
		t.Error(err)
	}

	if !db.Consents("logme") {
		t.Errorf("consent needs to be 1!!")
	}
}
