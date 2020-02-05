package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/xDarkicex/todos"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	"github.com/alecthomas/template"
	"github.com/go-chi/chi"
	"github.com/juju/errors"
	"github.com/xDarkicex/cchha_server_new/app/models"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

// Reviews type strictly for binding purposes
type Reviews Controllers

func init() {

}

type review struct {
	Featured         bool
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
	Details          []models.Detail
}

// Calc struct holds all calculation data
type Calc struct {
	Ratings []int
	Total
	Avg string
	Percentages
}

type CalcSorted struct {
	Ratings []int
	Total
	Avg string
	Percentages
}

// Total holds total star rating counts
type Total struct {
	Five  int
	Four  int
	Three int
	Two   int
	One   int
}

type Percentages struct {
	Five  string
	Four  string
	Three string
	Two   string
	One   string
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func Percent(total []int) map[int]float64 {
	var x []int
	var dict = make(map[int]int)
	for _, v := range total {
		x = append(x, v)
	}
	for _, num := range x {
		dict[num] = dict[num] + 1
	}
	var percentages = make(map[int]float64)
	var sum float64
	for k, v := range dict {
		fmt.Println(k, v)
		sum += float64(v)
		percentages[k] = toFixed((float64(v) / float64(len(total))), 2)
	}
	return percentages
}

func Sum(c []int) int {
	var temp int
	for i := 0; i < len(c); i++ {
		temp += c[i]
	}
	return temp
}

func sumPercise(c []int) float64 {
	var temp float64
	for i := 0; i < len(c); i++ {
		temp += float64(c[i])
	}
	return temp
}

// Mean is the average of all ratings
func Mean(c []int) int {
	var temp int
	temp = Sum(c)
	return temp / len(c)
}
func meanPercise(c []int) float64 {
	var sum float64
	sum = sumPercise(c)
	return toFixed(sum/float64(len(c)), 2)
}

// Index reviews list page
func (rev Reviews) Index(w http.ResponseWriter, r *http.Request) {
	var count = 0
	file, err := ioutil.ReadFile("./app/views/home-health/reviews.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, error.Error(errors.Trace(err)), http.StatusInternalServerError)
		return
	}

	data := models.GetApprovedReviews()
	if len(data) <= 0 {
		data = append(data, models.Review{Title: "No Reivews Found"})
		count = 0
	} else {
		count = len(data)
	}

	var temp = make([]int, 0)
	for _, post := range data {
		temp = append(temp, post.Rating)
	}
	var dict = make(map[int]int)
	for _, num := range temp {
		// Getting total count for each number set to key
		dict[num] = dict[num] + 1

	}
	p := message.NewPrinter(language.English)
	percentages := Percent(temp)
	c := Calc{
		Ratings: temp,
		Total: Total{
			One:   dict[1],
			Two:   dict[2],
			Three: dict[3],
			Four:  dict[4],
			Five:  dict[5],
		},
		Avg: p.Sprintf("%v", meanPercise(temp)),

		Percentages: Percentages{
			One:   p.Sprintf("%v", number.Percent(percentages[1])),
			Two:   p.Sprintf("%v", number.Percent(percentages[2])),
			Three: p.Sprintf("%v", number.Percent(percentages[3])),
			Four:  p.Sprintf("%v", number.Percent(percentages[4])),
			Five:  p.Sprintf("%v", number.Percent(percentages[5])),
		},
	}

	s, err := GetNamed(r, "current-session")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}

	user := helpers.GetCurrentUser("user-id", s)

	s.Save(r, w)
	t := template.Must(template.New("home-health").Funcs(funcMAP).Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Reviews",
		"Data":  data,
		"Calc":  c,
		"Count": count,
		"User":  user,
	})
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
}

func (rev Reviews) Sort(w http.ResponseWriter, r *http.Request) {
	var count = 0
	file, err := ioutil.ReadFile("./app/views/home-health/reviews.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, error.Error(errors.Trace(err)), http.StatusInternalServerError)
		return
	}
	s, _ := GetNamed(r, "current-session")
	numb := chi.URLParam(r, "sort")
	user := helpers.GetCurrentUser("user-id", s)
	r.ParseForm()
	rating, err := strconv.Atoi(numb)
	if err != nil {
		helpers.HandleError(err)
	}
	data := models.GetApprovedReviews()
	if len(data) <= 0 {
		data = append(data, models.Review{Title: "No Reivews Found"})
		count = 0
	} else {
		count = len(data)
	}
	var temp = make([]int, 0)
	for _, post := range data {
		temp = append(temp, post.Rating)
	}
	var dict = make(map[int]int)
	for _, num := range temp {
		// Getting total count for each number set to key
		dict[num] = dict[num] + 1

	}

	datasorted := models.GetApprovedReviewsSorted(rating)
	if len(datasorted) <= 0 {
		datasorted = append(datasorted, models.Review{Title: "No Reivews Found"})
	}
	var tempsorted = make([]int, 0)
	for _, post := range datasorted {
		tempsorted = append(tempsorted, post.Rating)
	}
	var dictsorted = make(map[int]int)
	for _, num := range tempsorted {
		// Getting total count for each number set to key
		dictsorted[num] = dictsorted[num] + 1

	}

	ps := message.NewPrinter(language.English)
	percentagesSorted := Percent(tempsorted)

	cs := CalcSorted{
		Ratings: tempsorted,
		Total: Total{
			One:   dictsorted[1],
			Two:   dictsorted[2],
			Three: dictsorted[3],
			Four:  dictsorted[4],
			Five:  dictsorted[5],
		},
		Avg: ps.Sprintf("%v", meanPercise(tempsorted)),

		Percentages: Percentages{
			One:   ps.Sprintf("%v", number.Percent(percentagesSorted[1])),
			Two:   ps.Sprintf("%v", number.Percent(percentagesSorted[2])),
			Three: ps.Sprintf("%v", number.Percent(percentagesSorted[3])),
			Four:  ps.Sprintf("%v", number.Percent(percentagesSorted[4])),
			Five:  ps.Sprintf("%v", number.Percent(percentagesSorted[5])),
		},
	}
	p := message.NewPrinter(language.English)
	percentages := Percent(temp)
	c := Calc{
		Ratings: temp,
		Total: Total{
			One:   dict[1],
			Two:   dict[2],
			Three: dict[3],
			Four:  dict[4],
			Five:  dict[5],
		},
		Avg: p.Sprintf("%v", meanPercise(temp)),

		Percentages: Percentages{
			One:   p.Sprintf("%v", number.Percent(percentages[1])),
			Two:   p.Sprintf("%v", number.Percent(percentages[2])),
			Three: p.Sprintf("%v", number.Percent(percentages[3])),
			Four:  p.Sprintf("%v", number.Percent(percentages[4])),
			Five:  p.Sprintf("%v", number.Percent(percentages[5])),
		},
	}

	s.Save(r, w)
	t := template.Must(template.New("home-health").Funcs(funcMAP).Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title":      "Home Health - Reviews",
		"Data":       datasorted,
		"Calc":       c,
		"CalcSorted": cs,
		"Count":      count,
		"User":       user,
	})
}

// Create control for crud opperations on reviews datatbase
func (rev Reviews) Create(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	user := helpers.GetCurrentUser("user-id", s)
	r.ParseForm()

	rating, err := strconv.Atoi(r.Form["rating"][len(r.Form["rating"])-1 : len(r.Form["rating"])][0])
	if err != nil {
		helpers.HandleError(err)
	}
	todo := elephant.GetElephant("todos")
	todo.SetMemory("Todo", "Build GuestReview struct and migrate without constrints")
	review := models.Review{
		Username:         r.Form.Get("review_name"),
		Email:            r.Form.Get("review_email"),
		Title:            r.Form.Get("review_title"),
		Body:             r.Form.Get("review_body"),
		Rating:           rating,
		UserID:           2,
		VisitorID:        user.ID,
		ExternalLink:     "https://Compassionatecare.com",
		ExternalSiteName: "Compassionate-Care-Home-Health",
		Pending:          true,
	}

	review = models.CreateReview(review)
	detail := models.Detail{
		ApprovalTime: time.Now(),
		ReviewID:     review.ID,
		UserID:       user.ID,
		ReviewerID:   user.ID,
		Title:        fmt.Sprintf("FROM_FRONT_END"),
		Body:         "New review from user",
	}
	detail = models.CreateDetail(detail)
	fmt.Println(detail)
	review.Details = append(review.Details, detail)
	models.UpdateReview(review)
	s.Save(r, w)
	http.Redirect(w, r, "/home-health/reviews", 302)
}

func (res Reviews) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "reviewID")
	sql := fmt.Sprintf("id= '%s'", id)
	p := models.GetReview(sql)
	d := models.GetReviewWithDetails(p.ID)
	data := struct {
		Review  models.Review
		Details []models.Detail
	}{
		Review:  p,
		Details: d,
	}
	encode := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := encode.Encode(data)
	if err != nil {
		helpers.HandleError(err)
	}
}

// Details ...
func (rev Reviews) Details(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "reviewID")
	idn, _ := strconv.Atoi(id)
	sql := fmt.Sprintf("id= '%d'", idn)
	p := models.GetReview(sql)
	d := models.GetReviewWithDetails(p.ID)
	data := struct {
		Details []models.Detail
	}{
		Details: d,
	}
	encode := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	encode.Encode(data)
}

func (rev Reviews) Json(w http.ResponseWriter, r *http.Request) {
	d := models.GetDetails()
	// u := models.GetUsers()
	p := models.GetReviews()

	data := struct {
		Reviews []models.Review
		Details []models.Detail
		// Users   []models.User
	}{
		Reviews: p,
		Details: d,
		// Users:   u,
	}
	encode := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	helpers.EnableCors(&w)
	// w.WriteHeader(http.StatusOK)
	err := encode.Encode(data)
	if err != nil {
		helpers.HandleError(err)
	}
}

func (reviews Reviews) Feature(w http.ResponseWriter, r *http.Request) {
	var request = struct {
		ReviewID uint `json:"review_id"`
		UserID   uint `json:"user_id"`
		Featured bool `json:"featured"`
	}{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	encoder := json.NewEncoder(w)
	err := decoder.Decode(&request)
	if err != nil {
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		err = encoder.Encode(&Res{Success: false, Message: "Malformed Request: " + errors.Cause(err).Error(), Time: time.Now().Local().Format(time.Stamp)})
		if err != nil {
			helpers.HandleError(err)
			return
		}
	}
	sql := fmt.Sprintf("id = %d", request.ReviewID)
	review := models.GetReview(sql)
	detail := models.Detail{
		ApprovalTime: time.Now(),
		ReviewID:     request.ReviewID,
		ReviewerID:   request.UserID,
		UserID:       request.UserID,
		Title:        "Change in status",
	}

	detail = models.CreateDetail(detail)
	review.Featured = request.Featured
	review.Details = append(review.Details, detail)
	err = models.UpdateReview(review)
	if err != nil {
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		err = encoder.Encode(&Res{Success: false, Message: "Malformed Request: " + errors.Cause(err).Error(), Time: time.Now().Local().Format(time.Stamp)})
		if err != nil {
			helpers.HandleError(err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(&Res{Success: true, Message: "Review set as Featured", Time: time.Now().Local().Format(time.Stamp)}); err != nil {
		helpers.HandleError(err)
	}
}
