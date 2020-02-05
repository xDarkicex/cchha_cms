package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	gs "github.com/gorilla/sessions"

	"github.com/xDarkicex/todos"

	"github.com/juju/errors"

	labstack "github.com/labstack/labstack-go"
	"github.com/labstack/labstack-go/ip"
	"github.com/xDarkicex/cchha_server_new/app/models"

	"github.com/alecthomas/template"
)

func HandleError(err error) {
	// notice that we're using 1, so it will actually log the where
	// the error happened, 0 = this function, we don't want that.
	pc, fn, line, _ := runtime.Caller(1)

	log.Println(Red("[error]"), errors.Cause(err))
	log.Printf(Red("[error]"), fmt.Sprintf(" in %s[%s:%d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
}

func SetCookie(name string, value interface{}, s *gs.Session) {
	s.Values[name] = value
}

func GetCurrentUser(name string, s *gs.Session) models.User {
	id, ok := s.Values[name].(uint)
	if !ok {
		HandleError(errors.NewForbidden(errors.New("unauthorized"), "Route Protected"))
		return models.User{}
	}
	spew.Dump(s.Values)
	fmt.Println(id)
	user, err := models.GetUser(fmt.Sprintf("id = %d", id))
	if err != nil {
		fmt.Println(err, "<<<", "\n[cause]", errors.Cause(err), "\n[Details]", errors.Details(err), "\n[stack]", errors.ErrorStack(err))
		// HandleError(err)
		return models.User{}
	}
	return user
}

func RedirectWithoutHTML(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" {
		return nil
	}
	if !strings.Contains(r.URL.Path, ".") {
		// force default all non-specific paths to .html
		http.Redirect(w, r, r.URL.Path+".html", http.StatusFound)
		return nil
	}
	return nil
}

func withoutHTML(w http.ResponseWriter, r *http.Request) string {
	// if r.URL.Path == "/" {
	// 	return nil
	// }
	// fmt.Println(r.URL.Path)
	if strings.Contains(r.URL.EscapedPath(), ".") {
		path := strings.Split(r.URL.EscapedPath(), ".")
		return path[0]
	}
	return r.URL.EscapedPath()
}

func Render(w http.ResponseWriter, r *http.Request) {
	path := withoutHTML(w, r)
	device := r.UserAgent()
	expression := regexp.MustCompile("(Mobi(le|/xyz)|Tablet)")
	if !expression.MatchString(device) {
		w.Header().Set("Connection", "keep-alive")
	}
	file, err := ioutil.ReadFile("./app/views/" + path + ".html")
	if err != nil {
		fmt.Println(err)
	}

	t, err := template.New(path).Parse(string(file))
	if err != nil {
		HandleError(err)
	}
	err = t.Execute(w, map[string]interface{}{})
	if err != nil {
		HandleError(err)
	}
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//HTTPS will redirect https traffic too http
func HTTPS(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.EscapedPath())
	target := "https://compassionatecare.com" + r.URL.EscapedPath()
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

// IsEmpty will return true is key is nil
func IsEmpty(key string) bool {
	// if key == "" {
	// 	return true
	// }
	// return false
	return key == ""
}

// labstack API KEY
// <_WS7Ro6_GOLrdzNhdKpUMcjwq4cWPfHrx_FgsAC6fevd25UB0EJt_>

func GetLocation(r *http.Request) models.Location {
	c := labstack.NewClient("_WS7Ro6_GOLrdzNhdKpUMcjwq4cWPfHrx_FgsAC6fevd25UB0EJt_")
	s := c.IP()
	res, err := s.Lookup(&ip.LookupRequest{
		IP: "104.70.66.209",
	})
	if err != nil {
		HandleError(errors.New(errors.Details(err)))
	}
	loc := models.Location{Lat: res.Latitude, Lng: res.Longitude, Country: res.Country, Time: res.TimeZone.Time, City: res.City, Region: res.Region, Postal: res.Postal}
	// TODOS CREATE MODEL LOCATION
	// models.CreateLocation(models.Location{Lat: res.Latitude, Lng: res.Longitude, Country: res.Country, Time: res.TimeZone.Time, City: res.City, Region: res.Region, Postal: res.Postal})
	return loc
}

func GetIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		HandleError(errors.Annotate(err, "error with host ip"))
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		log.Printf("%s", "remoteIP error net.Parse returned <nil>")
	}
	return remoteIP.String()
}

func Elephant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := elephant.GetElephant("todos")
		e.Emit()
		next.ServeHTTP(w, r)
	})
}
