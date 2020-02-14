package models

import (
	"log"
	"runtime"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type _review Review

func CreateReview(review Review) (rev Review) {
	review.ExternalSiteName = strings.ToLower(review.ExternalSiteName)
	err := db.Debug().Create(&review).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	db.Preload("Details").Find(&rev, review.ID)
	return rev
}

func (r Review) Delete() error {
	return db.Debug().Unscoped().Delete(&r).Error
}

func DeleteReview(sql string) error {
	review := Review{}
	return db.Debug().Delete(&review, sql).Error
}

func GetReviews() (reviews []Review) {
	db.Preload("Details").Find(&reviews)
	return reviews
}

func GetReview(sql string) (review Review) {
	// review := Review{}
	db.Preload("Details").Where(sql).Find(&review)
	return review
}

func GetReviewWithDetails(id uint) (details []Detail) {
	review := Review{}
	db.Preload("Details").Find(&review, id)
	return review.Details
}

func GetSeedReviews() (details []Detail) {
	review := Review{}
	db.Preload("Details").Find(&review, "user_id = 2")
	return review.Details
}

func GetReviewsByUser(id uint) (reviews []Review) {
	var reviews_fav = []Review{}
	var ordered_reviews = []Review{}
	err := db.Debug().Preload("Details").Where("user_id = $1 and featured = $2", id, true).Order("featured desc").Find(&reviews_fav).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	err = db.Debug().Preload("Details").Where("user_id = $1 and featured = $2", id, false).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	reviews = append(reviews_fav, ordered_reviews...)
	return reviews
}
func GetRejectedReviewsByUser(id uint) (reviews []Review) {
	var r []Review
	err := db.Preload("Details").Where("pending = $1 and user_id = $2", true, id).Order("rating desc").Find(&r).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	return r
}

func GetApprovedReviewsSorted(num int) (reviews []Review) {
	ordered_reviews := []Review{}
	ordered_reviews_2 := []Review{}

	//second in array
	err := db.Preload("Details").Where("rating = $1 AND pending = $2", num, false).Not("featured = ?", true).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	// first value
	err = db.Preload("Details").Where("featured = $1 AND rating = $2", true, num).Not("pending = $3", true).Order("rating desc").Find(&ordered_reviews_2).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	reviews = append(ordered_reviews_2, ordered_reviews...)
	return reviews
}

func GetPendingReviews() (reviews []Review) {
	reviews_fav := []Review{}
	ordered_reviews := []Review{}
	err := db.Preload("Details").Where("pending = ?", true).Not("featured = ?", true).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	err = db.Preload("Details").Where("featured = $1 AND pending= $2", true, true).Order("rating desc").Find(&reviews_fav).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	reviews = append(reviews_fav, ordered_reviews...)
	return reviews
}

func GetApprovedReviews() (reviews []Review) {
	reviews_fav := []Review{}
	ordered_reviews := []Review{}
	err := db.Preload("Details").Where("pending = ?", false).Not("featured = ?", true).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	err = db.Preload("Details").Where("featured = $1 AND pending= $2", true, false).Order("rating desc").Find(&reviews_fav).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	reviews = append(reviews_fav, ordered_reviews...)
	return reviews
}

func UpdateReview(review Review) error {
	err := db.Debug().Save(&review).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)

		return err
	}
	return nil
}
