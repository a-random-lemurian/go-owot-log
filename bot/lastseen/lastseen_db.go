package lastseen

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type LastseenUserEntry struct {
	gorm.Model
	RealUsername string
	Time         time.Time
	Consent      bool
}

type LastseenDB interface {
	GetLastseen(string) (time.Time, error)
	SetLastseen(string, time.Time) error
	SetConsent(string, bool) error
	Consents(string) bool
}

type GormLastseenDB struct {
	gorm *gorm.DB
}

func NewGormLastseenDB(gorm *gorm.DB) *GormLastseenDB {
	db := &GormLastseenDB{gorm: gorm}
	db.gorm.AutoMigrate(&LastseenUserEntry{})
	return db
}

var ErrNoConsent = errors.New("user did not consent to time lookups")

func (g *GormLastseenDB) GetLastseen(realUsername string) (time.Time, error) {
	consent := g.Consents(realUsername)
	if !consent {
		return time.Time{}, ErrNoConsent
	}

	query := LastseenUserEntry{RealUsername: realUsername}
	var user LastseenUserEntry
	tx := g.gorm.First(&user, query)
	if tx.Error != nil {
		return time.Time{}, tx.Error
	}
	return user.Time, nil
}

func (g *GormLastseenDB) SetLastseen(realUsername string, lastseenTime time.Time) error {
	var user LastseenUserEntry
	tx := g.gorm.Where(LastseenUserEntry{RealUsername: realUsername}).
		Assign(LastseenUserEntry{Time: lastseenTime}).
		FirstOrCreate(&user)
	return tx.Error
}

func (g *GormLastseenDB) SetConsent(realUsername string, consent bool) error {
	var user LastseenUserEntry
	tx := g.gorm.Where(LastseenUserEntry{RealUsername: realUsername}).
		Assign(LastseenUserEntry{Consent: consent}).
		FirstOrCreate(&user)
	return tx.Error
}

func (g *GormLastseenDB) Consents(realUsername string) bool {
	query := LastseenUserEntry{RealUsername: realUsername}
	var user LastseenUserEntry
	tx := g.gorm.First(&user, query)
	if tx.Error != nil {
		return true
	}
	return user.Consent
}
