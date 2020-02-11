package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	namer "github.com/xDarkicex/name_generator"

	"github.com/xDarkicex/name_generator/dictionary"

	"github.com/juju/errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xDarkicex/cchha_server_new/app/datastore"
)

var (
	SETUP_DB = false
	dbConn   *gorm.DB
	db       *datastore.Store
)

func init() {
	dbConn = datastore.Connect()
	db = datastore.New(dbConn)
	err := Migration(db)
	if err != nil {
		log.Println("[Cause] ", errors.Cause(err), "\n[Details] ", errors.Details(err))
	}
	// var review Review
	// db.Debug().Where("rating = 0").Unscoped().Delete(&review)
	// only run on seeding
	// should maybe add to flag
	// seed()
	// remove_seed()
}

func remove_seed() {
	time := time.Now()
	// DeleteDatail("user_id = 0")

	reviews := GetApprovedReviews()
	for _, review := range reviews {
		hour_db := review.Model.CreatedAt.Hour()
		day_db := review.Model.CreatedAt.Format("Monday")
		hour := time.Hour()
		day := time.Format("Monday")
		fmt.Println(day_db, " ", day, "\n", hour_db, " ", hour)
		if day_db == day {
			if hour_db == hour {
				for _, detail := range review.Details {
					if err := detail.Delete(); err != nil {
						fmt.Println(errors.Cause(err))
					}
				}
				if err := review.Delete(); err != nil {
					fmt.Println(errors.Cause(err))
				}
			}

		}
	}
}

func seed() {
	file := dictionary.RETURNSEED()
	// val := strings.Split(string(file), "},")
	var seeds = struct {
		Comments []struct {
			Name         string `json:"name,omitempty"`
			Title        string `json:"title,omitempty"`
			Body         string `json:"body"`
			Location     string `json:"location,omitempty"`
			Rating       int    `json:"rating"`
			ExternalSite string `json:"external_site"`
		} `json:"comments"`
	}{}
	err := json.Unmarshal(file, &seeds)
	if err != nil {
		fmt.Println("cause: ", errors.Cause(err), "\n", "Details:", errors.ErrorStack(err))
	}

	for _, seed := range seeds.Comments {

		review := Review{
			Rating:           seed.Rating,
			VisitorID:        0,
			Email:            namer.GetEmailDomain(),
			Title:            seed.Title,
			Body:             seed.Body,
			Username:         seed.Name,
			ExternalLink:     seed.Location,
			ExternalSiteName: seed.ExternalSite,
			UserID:           uint(0),
			Pending:          false,
		}
		// Details:          []Detail{detail},
		review = CreateReview(review)
		detail := Detail{
			ApprovalTime: time.Now(),
			UserID:       uint(0),
			ReviewerID:   uint(0),
			ReviewID:     review.ID,
			Title:        fmt.Sprintf("SEEDED"),
			Body:         fmt.Sprintf("Post Made from seed file"),
		}
		detail = CreateDetail(detail)
		review.Details = append(review.Details, detail)
		UpdateReview(review)

	}
}

////////////////////////////////////////////////////////////////////|
//																																	|
// Locations  																											|
// : Struct is for tracking location of visitors and users          |
//                                                                  |
//                                                                  |
/////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////|
//																																	|
// : Reviews  																						   				|
// : All data for reviews and the details about them                |
//                                                                  |
//                                                                  |
/////////////////////////////////////////////////////////////////////

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
	UserID           uint
	VisitorID        uint
	Details          []Detail `gorm:"many2many:review_details;"`
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

////////////////////////////////////////////////////////////////////|
//																																	|
// PageViews  																											|
// : To log and agregate page view, time and locaion in local env   |
// : allows for real time dashboard stats                           |
//                                                                  |
/////////////////////////////////////////////////////////////////////

type Page struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Title       string
	Path        string
	Screenshoot string
}

type View struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Unique    bool
	VisitorID uint
	StartTime time.Time
	EndTime   time.Duration
}

// Visitor ...
type Visitor struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	IP        string
	Returning bool
	Location
	Device string
	Pages  []Page `gorm:"many2many:visitor_pages;"`
	Date   time.Time
	Time   time.Duration
}

////////////////////////////////////////////////////////////////////|
//																																	|
// Flash  			     																								|
// : For sending contextual messages to the front end client        |
// : Used with Sessions                                             |
//                                                                  |
/////////////////////////////////////////////////////////////////////

type Flash struct {
	Type    string
	Message string
}

////////////////////////////////////////////////////////////////////|
//																																	|
// Files and resources			   																			|
// : verifying files and locations on server on upload and download |
// : Uses:URI, Hash check, Upload/Download, get location            |
//                                                                  |
/////////////////////////////////////////////////////////////////////

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
	File       `gorm:"embedded:"`
}

type Photo struct {
	gorm.Model
	UserPhoto uint
	Small     string
	Medium    string
	Large     string
	File
}

////////////////////////////////////////////////////////////////////|
//																																	|
// File and resource  			   																			|
// : verifying files and locations on server on upload and download |
// : Uses:URI, Hash check, Upload/Download, get location            |
//                                                                  |
/////////////////////////////////////////////////////////////////////

type Message struct {
	gorm.Model
	Title      string
	Type       uint // Has four states: Note, Notification, Direct Message, Group Message
	Body       string
	UserID     uint
	SenderID   uint
	RecevierID uint
	SentTime   time.Time
	Status     uint // four states: Sent, Unread, Read, Replied
}

////////////////////////////////////////////////////////////////////|
//																																	|
// Reset  			   																			            |
// : Reset is for reset token for password                          |
// :                                                                |
//                                                                  |
/////////////////////////////////////////////////////////////////////

type Reset struct {
	gorm.Model
	Token   string
	UserID  uint
	Active  bool
	Expires time.Time
}

////////////////////////////////////////////////////////////////////|
//																																	|
// Users  			   																			            |
// : All user data                                                  |
// :                                                                |
//                                                                  |
/////////////////////////////////////////////////////////////////////

type ACareer struct {
	Title string
	Years int
	Bio   string
}
type Avatar struct {
	Description string
	Photo       Photo
}
type AProfile struct {
	Avatar
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Age       int
	Phone     string
	Country   string
	Language  string
	Zip       string
	State     string
	City      string
	Street    string
	Photos    []Photo `gorm:"many2many:profile_photos;"`
}
type ASecurity struct {
	Email         string   `gorm:"unique;not null"`
	Password      string   `gorm:"not null"`
	IsSuperUser   bool     `gorm:"default:false"`
	IsAdmin       bool     `gorm:"default:false"`
	IsEditor      bool     `gorm:"default:false"`
	Location      Location `gorm:"type:json"`
	LastLogin     time.Time
	LastIP        string
	LoginAttempts int
}
type AUser struct {
	gorm.Model
	ACareer
	AProfile
	Security ASecurity
	Reviews  []Review  `gorm:"many2many:user_reviews;"`
	Messages []Message `gorm:"many2many:user_messages;"`
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
}
