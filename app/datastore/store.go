package datastore

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// var db = &Store{}

type Store struct {
	*gorm.DB
}

func New(conn *gorm.DB) *Store {
	return &Store{conn}
}

func Connect() *gorm.DB {

	db, err := gorm.Open("postgres", "host=localhost port=5432  sslmode=disable user=gentryrolofson dbname=cchha password=")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
