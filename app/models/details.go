package models

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type _Detail Detail

// GetDetail ...
func GetDetail(sql string) (deet Detail) {
	db.Find(&deet, sql)
	fmt.Println(deet, "<< get detail")
	return deet
}

func CreateDetail(detail Detail) (deet Detail) {
	db.Create(&detail)
	db.Find(&deet, "id = $1", detail.ID)
	return deet
}

// GetDetails Bd function !
func GetDetails() (details []Detail) {
	db.Find(&details)
	return details
}
