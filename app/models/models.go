package models

import (
	"log"
	"time"

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
