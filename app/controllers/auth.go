package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juju/errors"

	"github.com/go-chi/chi"
	"github.com/xDarkicex/cchha_server_new/app/models"
)

// Auth binding type
type Auth Controllers

// Authtenicate middleware
func (a Auth) Authtenicate(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	// everything here happpens on server Load
	var allUsers = []models.User{}
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Context().Value("requestID"))
		//Not working because the ID doesnt exist on all routes.
		if len(allUsers) == 0 && r.RequestURI == "/bd-admin/signup" {
			fn(w, r)
			return
		}
		ID := chi.URLParam(r, "userID")
		s, _ := GetNamed(r, "current-session")
		f, _ := GetNamed(r, "Flash")
		current_user_id, ok := s.Values["user-id"].(uint)
		if !ok {
			AddFlash(r, w, models.Flash{Type: "Warn", Message: "Not logged in."}, f)
			http.Redirect(w, r, "/", 302)
			return
		}
		user, err := models.GetUser(fmt.Sprintf("id = %d", current_user_id))
		if err != nil {
			http.Error(w, errors.Details(err), http.StatusInternalServerError)
			return
		}

		if ID != "" {
			fmt.Printf("User ID:%s\n", ID)
			id, err := strconv.Atoi(ID)
			if err != nil {
				http.Error(w, errors.Details(err), http.StatusInternalServerError)
				return
			}

			fmt.Println(user, "<<< user struct inside")
			// check for super user status
			if user.IsSuperUser == true {
				fmt.Println("super redirect")
				fn(w, r)
				return
				// check if user has rights to page access
			}
			if user.ID == uint(id) {
				fmt.Println("ids match redirect")
				fn(w, r)
				return
			}
		}

		// f := fmt.Sprintf("Hello: %s, Super Status: %v\nAdmin Status: %v\n", user.FirstName, user.IsSuperUser, user.IsAdmin)
		// fmt.Print(f)
		if user.IsSuperUser == true {
			// fmt.Printf("super")
			fn(w, r)
			return
		}
		if user.IsAdmin == true {
			// fmt.Println("Is admin User")
			fn(w, r)
			return
		}
		http.Redirect(w, r, "/bd-admin/signin", 302)
		return
	}
}
