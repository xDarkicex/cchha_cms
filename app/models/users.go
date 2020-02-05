package models

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/errors"
)

type _user User

func CreateUser(user User) error {
	err := db.Debug().Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Update ...
func UpdateUser(user User) {
	db.Save(user)
}

func GetUser(sql string) (u User, err error) {
	err = db.Where(sql).Find(&u).Error
	if err != nil {
		fmt.Println("======", errors.Details(err), "======")
		return User{}, err
	}
	return u, nil
}

func GetUsers() {

}
