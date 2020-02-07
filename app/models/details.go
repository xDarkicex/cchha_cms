package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type _Detail Detail

// GetDetail ...
func GetDetail(sql string) (deet Detail) {
	db.Debug().Find(&deet, sql)
	return deet
}

func CreateDetail(detail Detail) (det Detail) {
	db.Debug().Create(&detail)
	db.Debug().Find(&det, "id = $1", detail.ID)
	return det
}

// GetDetails Bd..
func GetDetails() (details []Detail) {
	db.Debug().Find(&details)
	return details
}

func (d Detail) Delete() error {
	return db.Debug().Unscoped().Delete(&d).Error
}

func DeleteDatail(sql string) error {
	detail := Detail{}
	return db.Debug().Delete(&detail, sql).Error
}
