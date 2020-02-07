package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/alecthomas/template"
	"github.com/juju/errors"
	"github.com/xDarkicex/cchha_server_new/app/models"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

// change rating post time format
const customFormat = `January _2, 2006`

var funcMAP = template.FuncMap{
	"GetFlashes": func(flashes []models.Flash) []models.Flash {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		for _, k := range flashes {
			var flash models.Flash
			enc.Encode(k)
			json.Unmarshal(buf.Bytes(), &flash)
			flashes = append(flashes, flash)
			buf.Reset()
		}
		return flashes
	},
	"GetNewestFlashes": func(flashes []models.Flash) map[string]models.Flash {
		var MostRecent = make(map[string]models.Flash)
		for k, v := range flashes {

			if len(flashes)-1 == k {
				MostRecent[v.Type] = flashes[k]
			}
			if len(flashes)-2 == k {
				MostRecent[v.Type] = flashes[k]
			}
			if len(flashes)-3 == k {
				MostRecent[v.Type] = flashes[k]
			}
		}
		return MostRecent
	},
	"ToDate": func(t time.Time) string {
		return t.Format(customFormat)
	},
	"ToTime": func(t time.Time) string {
		return t.Format("3:04PM")
	},
	"ToString": func(approved []models.Review) string {
		var data []string
		var temp []byte
		for i := range approved {
			temp, _ = json.Marshal(&approved[i])
			data = append(data, string(temp))
		}
		reviews, _ := json.Marshal(data)
		return string(reviews)
	},
}

type Application Controllers

func (a Application) Index(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/splash.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}

	t := template.Must(template.New("index").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Title": "Compassionate Care - Splash",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}

func (a Application) CustomNotFound(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("./app/views/404.html")
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
	// w.WriteHeader(http.StatusNotFound)
	t := template.Must(template.New("error_page").Parse(string(file)))
	err = t.Execute(w, map[string]interface{}{
		"Error": r.Host + r.URL.EscapedPath() + " page not found",
	})
	if err != nil {
		helpers.HandleError(err)
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
	}
}
