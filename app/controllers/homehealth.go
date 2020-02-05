package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/juju/errors"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

type HomeHealth Controllers

func (this HomeHealth) Index(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/index.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home health - Home Page",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Careers(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/careers.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Careers",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Services(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/services.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Services",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Eligibility(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/eligibility.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Eligibility",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Community(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/community.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Community",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Resources(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/resources.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Resources",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Contact(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/contact.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Contact",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) Locations(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/locations.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - Locations",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (this HomeHealth) About(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/home-health/about.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("home-health").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Home Health - About",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}
