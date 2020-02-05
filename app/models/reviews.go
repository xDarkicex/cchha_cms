package models

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type _review Review

func CreateReview(review Review) (rev Review) {
	review.ExternalSiteName = strings.ToLower(review.ExternalSiteName)
	fmt.Println(review)
	err := db.Debug().Create(&review).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	db.Preload("Details").Find(&rev, review.ID)
	return rev
}

func DeleteReview(review Review) error {
	return db.Delete(&review).Error
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

func GetApprovedReviewsByUser(id uint) (reviews []Review) {
	err := db.Preload("Details").Where("pending = ?", false).Where("user_id = ?", id).Order("rating desc").Find(&reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	return reviews
}

// , id, db.Select("*").Table("reviews").Where("pending = ?", true).QueryExpr()).Row()
// 	results := []Review{}
// 	err := row.Scan(&results).Error
// db.Debug().Where("pending = $1", true).Table("reviews").Find(&r)
// 	db.Debug().Select("user_id, review_id, reviewer_id").Where("user_id = $1", id).Find(&result)
// db.Debug().Table("reviews").Preload("Details").Where("user_id = $1 and reviews.pending = $2", id, true).Find(&r)
// db.Debug().Table("review_details").Where("")
// if err != nil {
// 	pc, fn, line, _ := runtime.Caller(1)
// 	log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
// // }
// err := db.Table("reviews").Preload("Details").Where("pending = ?", true).Where("details.user_id = ?", id).Order("rating desc").Find(&reviews).Error
func GetRejectedReviewsByUser(id uint) (reviews []Review) {
	// result := []Detail{}
	// r := []Review{}

	return nil
	// return reviews
}

func GetApprovedReviewsSorted(num int) (reviews []Review) {
	ordered_reviews := []Review{}
	err := db.Preload("Details").Where("rating = $1", num).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	return ordered_reviews
}

func GetPendingReviews() (reviews []Review) {
	pending_reviews := []Review{}
	err := db.Preload("Details").Where("pending = ?", true).Order("rating desc").Find(&pending_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	spew.Dump(pending_reviews)
	return pending_reviews
}

func GetApprovedReviews() (reviews []Review) {
	reviews_fav := []Review{}
	ordered_reviews := []Review{}
	err := db.Preload("Details").Where("pending = ?", false).Not("featured = ?", true).Order("rating desc").Find(&ordered_reviews).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	err = db.Preload("Details").Where("featured = ?", true).Order("rating desc").Find(&reviews_fav).Error
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	reviews = append(reviews_fav, ordered_reviews...)
	return reviews
}

func ApproveReview() {

}

// RejectReview -- Good function
func RejectReview() {

}

func UpdateReview(review Review) error {
	err := db.Save(&review).Error
	if err != nil {
		fmt.Println("<< Update Review")
		fmt.Println()
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)

		return err
	}
	return nil
}
