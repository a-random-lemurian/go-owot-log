package lastseen

import (
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

func (g *GormLastseenDB) GetLastseen(realUsername string) (time.Time, error) {
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
		Attrs(LastseenUserEntry{Consent: true}).
		FirstOrInit(&user)

	if tx.Error != nil {
		return tx.Error
	}

	user.Time = lastseenTime
	return g.gorm.Save(&user).Error
}

func (g *GormLastseenDB) SetConsent(realUsername string, consent bool) error {
	var user LastseenUserEntry
	tx := g.gorm.Where(LastseenUserEntry{RealUsername: realUsername}).First(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			user = LastseenUserEntry{
				RealUsername: realUsername,
				Consent:      consent,
			}
			return g.gorm.Create(&user).Error
		}
		return tx.Error
	}

	user.Consent = consent
	return g.gorm.Save(&user).Error
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
