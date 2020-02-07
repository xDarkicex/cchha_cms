package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/alecthomas/template"
	"github.com/go-chi/chi"
	"github.com/xDarkicex/cchha_server_new/app/models"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

// Admin ...
type Admin Controllers

// Index shows admin dashboard
func (a Admin) Index(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	f, _ := GetNamed(r, "Flash")

	var count = 0
	// spew.Dump(s.Values)
	user := helpers.GetCurrentUser("user-id", s)

	file, err := ioutil.ReadFile("./app/views/admin/Dashboard.html")
	if err != nil {
		helpers.HandleError(err)
	}
	data := models.GetPendingReviews()
	fmt.Println(data)
	if len(data) <= 0 {
		data = append(data, models.Review{Title: "No Reviews Found"})
		count = 0
	} else {
		count = len(data)
	}
	flashes := GetFlashes(r, f)
	err = f.Save(r, w)
	if err != nil {
		helpers.HandleError(err)
	}
	s.Save(r, w)
	t := template.Must(template.New("admin").Funcs(funcMAP).Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title":   "Admin Dashboard",
		"Data":    data,
		"User":    user,
		"Count":   count,
		"Flashes": flashes,
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (a Admin) Edit(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	userID := chi.URLParam(r, "userID")
	sql := fmt.Sprintf("id = %s", userID)
	user, err := models.GetUser(sql)
	if err != nil {
		helpers.HandleError(errors.Wrap(err, errors.New("Will go away once user tables have been migrated")))
	}

	online_user := helpers.GetCurrentUser("user-id", s)
	if online_user.IsSuperUser != true {
		if user.ID != online_user.ID {
			http.Error(w, "403: Forbidden Access", http.StatusForbidden)
			return
		}
	}
	if r.Method == "GET" {
		file, err := ioutil.ReadFile("./app/views/admin/edit.html")
		if err != nil {
			helpers.HandleError(err)
			http.Error(w, errors.Details(err), http.StatusInternalServerError)
		}
		s.Save(r, w)
		t := template.Must(template.New("admin").Funcs(funcMAP).Parse(string(file)))
		err = t.Execute(w, map[string]interface{}{
			"Title":    "Admin Edit User Data",
			"User":     user,
			"Loggedin": online_user,
		})
		if err != nil {
			helpers.HandleError(err)
			http.Error(w, errors.Details(err), http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println(r.Form)
		if !helpers.IsEmpty(r.Form.Get("first_name")) {
			user.FirstName = r.Form.Get("first_name")
		}
		if !helpers.IsEmpty(r.Form.Get("last_name")) {
			user.LastName = r.Form.Get("last_name")
		}
		if !helpers.IsEmpty(r.Form.Get("email")) {
			user.Email = r.Form.Get("email")
		}
		if !helpers.IsEmpty(r.Form.Get("phone")) {
			user.Phone = r.Form.Get("phone")
		}
		if !helpers.IsEmpty(r.Form.Get("title")) {
			user.Title = r.Form.Get("title")
		}
		if !helpers.IsEmpty(r.Form.Get("password")) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Form.Get("password")), bcrypt.DefaultCost)
			if err != nil {
				helpers.HandleError(err)
			}
			user.Password = string(hashedPassword)
		}
		if !helpers.IsEmpty(r.Form.Get("admin")) {
			user.IsAdmin = true
		}
		if helpers.IsEmpty(r.Form.Get("admin")) {
			user.IsAdmin = false
		}
		if !helpers.IsEmpty(r.Form.Get("editor")) {
			user.IsEditor = true
		}
		if helpers.IsEmpty(r.Form.Get("editor")) {
			user.IsEditor = false
		}
		s.Values["user-id"] = user.ID
		models.UpdateUser(user)
		err := s.Save(r, w)
		if err != nil {
			helpers.HandleError(err)
		}
		url := fmt.Sprintf("/bd-admin/user/%d", user.Model.ID)
		http.Redirect(w, r, url, 302)
		return
	}
}

// Show ...
func (a Admin) Show(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	// userID := chi.URLParam(r, "userID")
	// sql := fmt.Sprintf("id = '%s'", userID)
	// exists, err := models.GetUser(sql)
	// if err != nil {
	// 	helpers.HandleError(errors.New(errors.Details(err)))
	// }
	user := helpers.GetCurrentUser("user-id", s)
	// if exists.ID != user.ID {
	// 	http.Redirect(w, r, "/", 302)
	// 	return
	// }

	fmt.Println(user.ID, "<< in admin controller ===")
	rejected := models.GetRejectedReviewsByUser(user.ID)
	reviews := models.GetApprovedReviewsByUser(user.ID)

	file, err := ioutil.ReadFile("./app/views/admin/show-user.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("admin/show").Funcs(funcMAP).Parse(string(file)))
	fmt.Println(reviews)
	fmt.Println(rejected)
	details := models.GetDetails()
	s.Save(r, w)
	err = t.Execute(w, map[string]interface{}{
		"Title":    "Admin Profile",
		"User":     user,
		"Approved": reviews,
		"Rejected": rejected,
		"Notes":    details,
		"Count":    len(reviews),
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

// Create ...
func (a Admin) Create(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	user := helpers.GetCurrentUser("user-id", s)
	r.ParseForm()
	rawRating := r.Form["rating"]
	rating, err := strconv.Atoi(rawRating[len(rawRating)-1])

	if err != nil {
		helpers.HandleError(err)
	}
	review := models.Review{
		Rating:           rating,
		VisitorID:        0,
		Email:            r.Form.Get("email"),
		Title:            r.Form.Get("review_title"),
		Body:             r.Form.Get("review_body"),
		Username:         r.Form.Get("review_name"),
		ExternalLink:     r.Form.Get("external_site_link"),
		ExternalSiteName: r.Form.Get("external_site_name"),
		UserID:           user.ID,
		Pending:          false,
	}
	// Details:          []models.Detail{detail},
	review = models.CreateReview(review)
	detail := models.Detail{
		ApprovalTime: time.Now(),
		UserID:       user.ID,
		ReviewerID:   user.ID,
		ReviewID:     review.ID,
		Title:        fmt.Sprintf("MANUALLY_CREATED"),
		Body:         fmt.Sprintf("Post Made From Admin Panel by %s %s", user.FirstName, user.LastName),
	}
	detail = models.CreateDetail(detail)
	review.Details = append(review.Details, detail)
	models.UpdateReview(review)
	s.Save(r, w)
	http.Redirect(w, r, "/bd-admin", 302)
}

/////////////////////////////////|
//                               |
//       <= API CALLS >=         |
//                               |
//   THE FUTURE IS INTERACTIVE   |
/////////////////////////////////|

type Command struct {
	ReviewID int    `json:"review_id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

type Res struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

var Detail models.Detail

func (a Admin) Reject(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	f, _ := GetNamed(r, "Flash")
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	var msg = Command{}
	err := decoder.Decode(&msg)
	if err != nil {
		AddFlash(r, w, models.Flash{Type: "Warn", Message: fmt.Sprintf("Malformed Request: %s", errors.Cause(err))}, f)
		helpers.HandleError(err)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		err = f.Save(r, w)
		if err != nil {
			helpers.HandleError(err)
		}
		w.WriteHeader(http.StatusTeapot)
		err = encoder.Encode(&Res{Success: false, Message: fmt.Sprintf("Malformed Request: %s", errors.Details(err)), Time: time.Now().Local().Format(time.Stamp)})
		if err != nil {
			helpers.HandleError(err)
		}
	}
	user := helpers.GetCurrentUser("user-id", s)
	spew.Dump(user)
	if err != nil {
		helpers.HandleError(err)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		err = encoder.Encode(&Res{Success: false, Message: fmt.Sprintf("Malformed Request: %s", errors.Details(err)), Time: time.Now().Local().Format(time.Stamp)})
		if err != nil {
			helpers.HandleError(err)
		}
	}
	n := models.Detail{
		RejectionTime: time.Now(),
		ReviewID:      uint(msg.ReviewID),
		UserID:        user.ID,
		ReviewerID:    user.ID,
		Body:          msg.Body,
		Title:         msg.Title,
	}
	detail := models.CreateDetail(n)
	review := models.GetReview(fmt.Sprintf("id = '%d'", msg.ReviewID))
	reviews := models.GetReviews()
	for _, rev := range reviews {
		if rev.ID == review.ID {
			fmt.Println(rev.ID, review.ID)
			review.Pending = true
			fmt.Println(len(review.Details))
			review.Details = append(review.Details, detail)
			fmt.Println(len(review.Details))
		}
	}
	err = models.UpdateReview(review)
	if err != nil {
		helpers.HandleError(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = encoder.Encode(&Res{Success: true, Message: fmt.Sprintf("Rejection Accepted"), Time: time.Now().Local().Format(time.Stamp)})
	if err != nil {
		helpers.HandleError(err)
	}
}

// Approve serves as create route
func (a Admin) Approve(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	var msg = Command{}
	println(r.FormValue("title"))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	encoder := json.NewEncoder(w)
	err := decoder.Decode(&msg)
	if err != nil {
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		err = encoder.Encode(&Res{Success: false, Message: "Malformed Request: " + errors.Cause(err).Error(), Time: time.Now().Local().Format(time.Stamp)})
		if err != nil {
			helpers.HandleError(err)
		}
		http.Redirect(w, r, "/bd-admin", 300)
		return
	}
	user := helpers.GetCurrentUser("user-id", s)

	n := models.Detail{
		ApprovalTime: time.Now(),
		UserID:       user.ID,
		ReviewerID:   user.ID,
		ReviewID:     uint(msg.ReviewID),
		Title:        msg.Title,
		Body:         msg.Body,
	}

	detail := models.CreateDetail(n)
	review := models.GetReview(fmt.Sprintf("id = '%d'", msg.ReviewID))
	reviews := models.GetReviews()
	for _, rev := range reviews {
		if rev.ID == review.ID {
			fmt.Println(rev.ID, review.ID)
			review.Pending = false
			fmt.Println(len(review.Details))
			review.Details = append(review.Details, detail)
			fmt.Println(len(review.Details))
		}
	}

	err = models.UpdateReview(review)
	if err != nil {
		helpers.HandleError(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(&Res{Success: true, Message: "Approval Accepted", Time: time.Now().Local().Format(time.Stamp)}); err != nil {
		helpers.HandleError(err)
	}
}

func untitled(title string) string {
	if title != "" {
		return title
	}
	return ""
}

func (a Admin) DeleteReview(w http.ResponseWriter, r *http.Request) {
	// s, _ := GetNamed(r, "current-session")
	f, _ := GetNamed(r, "Flash")
	var msg = Command{}
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	err := decoder.Decode(&msg)
	if err != nil {
		AddFlash(r, w, models.Flash{Type: "Warn", Message: fmt.Sprintf("Malformed Request: %s", errors.Details(err))}, f)
		helpers.HandleError(err)
	}
	review := models.GetReview(fmt.Sprintf("id = '%d'", msg.ReviewID))
	err = review.Delete()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(&Res{Success: true, Message: "Approval Accepted", Time: time.Now().Local().Format(time.Stamp)}); err != nil {
		helpers.HandleError(err)
	}
}
