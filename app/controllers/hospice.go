package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"

	"github.com/alecthomas/template"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

type Hospice Controllers

func (this Hospice) Index(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/index.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - home Page",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Careers(w http.ResponseWriter, r *http.Request) {
	// helpers.RedirectWithoutHTML(w, r)
	// helpers.Render(w, r, "hospice/careers")
	file, err := ioutil.ReadFile("./app/views/hospice/careers.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Hospice - Careers",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Services(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/services.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Services",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Eligibility(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/eligibility.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Eligibility",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Community(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/community.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Community",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Resources(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/resources.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Resources",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Contact(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/contact.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Contact",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) Locations(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/locations.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "hospice - Locations",
	})
	if err != nil {
		helpers.HandleError(err)
	}
}

func (this Hospice) About(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/hospice/about.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	t := template.Must(template.New("hospice").Parse(string(file)))
	t.Execute(w, map[string]interface{}{
		"Title": "hospice - About",
	})
}
