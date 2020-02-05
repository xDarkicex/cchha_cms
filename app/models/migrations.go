package models

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/xDarkicex/cchha_server_new/app/datastore"

	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

func Migration(db *datastore.Store) error {
	m := gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "New Database - 2",
			Migrate: func(tx *gorm.DB) error {
				type Location struct {
					gorm.Model
					Lat     float64
					Lng     float64
					Country string
					Region  string
					City    string
					Postal  string
					Time    time.Time
				}
				type File struct {
					Name      string
					Ext       string
					URI       string
					Directory string
					Size      int64
					ModeTime  time.Time
					Temporary bool   `gorm:"default:false;"`
					Hash      string `gorm:"unique;not null"`
				}
				type Upload struct {
					Identity   int `gorm:"primary_key;"`
					UploaderID int
					Time       time.Time
					Approved   bool `gorm:"default:true;"`
					File
				}
				type Photo struct {
					gorm.Model
					Small  string
					Medium string
					Large  string
					File
				}
				type Avatar struct {
					Description string
					UserID      uint
					Photo
				}
				type Detail struct {
					gorm.Model
					ApprovalTime  time.Time
					RejectionTime time.Time
					ReviewID      uint
					ReviewerID    uint
					UserID        uint
					Title         string
					Body          string
				}
				type Review struct {
					gorm.Model
					Title            string
					Rating           int
					Username         string
					Email            string
					Body             string
					ExternalLink     string
					ExternalSiteName string
					Pending          bool
					UserID           uint
					VisitorID        uint
					Details          []Detail `gorm:"many2many:review_details;"`
				}
				type User struct {
					gorm.Model
					Title         string
					Years         int
					Bio           string
					FirstName     string `gorm:"not null"`
					LastName      string `gorm:"not null"`
					Age           int
					Phone         string
					Country       string
					Language      string
					Zip           string
					State         string
					City          string
					Street        string
					Description   string
					Email         string   `gorm:"unique;not null"`
					Password      string   `gorm:"not null"`
					IsSuperUser   bool     `gorm:"default:false"`
					IsAdmin       bool     `gorm:"default:false"`
					IsEditor      bool     `gorm:"default:false"`
					Location      Location `gorm:"type:json"`
					LastLogin     time.Time
					LastIP        string
					LoginAttempts int
					// Reviews       []Review  `gorm:"many2many:user_reviews;"`
					// Messages      []Message `gorm:"many2many:user_messages;"`
				}
				type Message struct {
					gorm.Model
					Title      string
					Type       uint
					Body       string
					SenderID   uint
					RecevierID uint
					SentTime   time.Time
					Status     uint
				}
				err := tx.Debug().CreateTable(&User{}, &File{}, &Upload{}, &Photo{}, &Avatar{}, &Message{},
					&Location{}, &Detail{}, &Review{}).Error
				if err != nil {
					fmt.Println("====", errors.Trace(err), "==== line 295")
				}
				db.Model(&Avatar{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
				return db.DB.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("careers", "avatars", "profiles", "securities", "users").Error
			},
		},
		{
			ID: "relations:users:photo:message:reivews",
			Migrate: func(tx *gorm.DB) error {
				type Photo struct {
					UserPhoto uint
				}
				type Message struct {
					UserID uint
				}
				err := tx.Debug().AutoMigrate(&Photo{}, &Message{}).Error
				if err != nil {
					fmt.Println("====", errors.Trace(err), "==== line 295")
				}
				db.Model(&Photo{}).AddForeignKey("user_photo", "users(id)", "CASCADE", "CASCADE")
				db.Model(&Message{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
				return db.DB.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("careers", "avatars", "profiles", "securities", "users").Error
			},
		},
		{
			ID: "relations:reviews:detail 8",
			Migrate: func(tx *gorm.DB) error {
				type Detail struct {
					gorm.Model
					ApprovalTime  time.Time
					RejectionTime time.Time
					ReviewID      uint `gorm:"unique;"`
					ReviewerID    uint
					UserID        uint
					Title         string
					Body          string
				}
				type Review struct {
					gorm.Model
					Featured         bool `gorm:"default:false;"`
					Title            string
					Rating           int
					Username         string
					Email            string
					Body             string
					ExternalLink     string
					ExternalSiteName string
					Pending          bool
					UserID           uint `gorm:"unique;"`
					VisitorID        uint
					Details          []Detail `gorm:"many2many:review_details;"`
				}

				err := tx.Debug().AutoMigrate(&Detail{}, &Review{}).Error
				if err != nil {
					fmt.Println("====", errors.Trace(err), "==== line 295")
				}
				db.Debug().Model(&Detail{}).AddForeignKey("review_id", "reviews(id)", "CASCADE", "CASCADE")
				return db.DB.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("careers", "avatars", "profiles", "securities", "users").Error
			},
		},
		{
			ID: "Add Featured Review",
			Migrate: func(tx *gorm.DB) error {
				type Review struct {
					Featured bool `gorm:"default:false;"`
				}
				err := tx.Debug().AutoMigrate(&Review{}).Error
				if err != nil {
					fmt.Println("====", errors.Cause(err), "====")
				}
				return db.DB.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("reviews").Error
			},
		},
		{
			ID: "Add keys 4",
			Migrate: func(tx *gorm.DB) error {
				type Review struct {
					gorm.Model
					Title            string
					Rating           int
					Username         string
					Email            string
					Body             string
					ExternalLink     string
					ExternalSiteName string
					Pending          bool
					UserID           uint
					VisitorID        uint
					Featured         bool     `gorm:"default:false;"`
					Details          []Detail `gorm:"many2many:review_details;"`
				}
				err := tx.Debug().AutoMigrate(&Review{}).Error
				if err != nil {
					fmt.Println("====", errors.Trace(err), "==== line 295")
				}
				db.Model(&Detail{}).AddForeignKey("user_id", "details(id)", "CASCADE", "CASCADE")
				db.Model(&Review{}).AddForeignKey("user_id", "reviews(id)", "CASCADE", "CASCADE")
				return db.DB.Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("reviews", "details").Error
			},
		},
	})
	return m.Migrate()
}
