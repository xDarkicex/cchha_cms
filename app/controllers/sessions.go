package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/template"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/securecookie"
	gs "github.com/gorilla/sessions"

	"github.com/xDarkicex/cchha_server_new/app/models"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

type Sessions Controllers

// KEY is global key for sessions
var KEY []byte

type Config struct {
	KEY string `json:"key"`
}

var _store *gs.CookieStore

func init() {
	//TODOS
	fmt.Println(helpers.Red("Config.ini"), "make config ini file == controllers/sessions.go line 32")
	// cfg, err := ini.Load("config/config.ini")
	// if err != nil {
	// 	fmt.Printf("Error loading ini: %s\n\r", err)
	// 	os.Exit(1)
	// }
	KEY = getKey()
	_store = Store()
	// if the file exist dont overwrite

	var c = Config{}
	err := json.Unmarshal(KEY, &c)
	if err != nil {
		helpers.HandleError(err)
	}
	fmt.Println(c, "=== key ===")
}

// KEY dont lose yours keys
func getKey() []byte {
	exists, err := os.Stat("key.json")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	fmt.Println(exists.Name())
	if exists != nil {
		f, err := ioutil.ReadFile("key.json")
		if err != nil {
			helpers.HandleError(errors.Cause(err))
		}
		return f
	}
	k := securecookie.GenerateRandomKey(32)
	key := Config{
		KEY: string(k),
	}
	data, err := json.Marshal(key)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	err = ioutil.WriteFile("key.json", data, 0644)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	jsonData, err := ioutil.ReadFile("key.json")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	return jsonData

}

func GetNamed(req *http.Request, name string) (*gs.Session, error) {
	return _store.Get(req, name)
}

func GetCookie(name string, s *gs.Session) interface{} {
	val := s.Values[name]
	return val
}

func SetCookie(name string, value interface{}, s *gs.Session) {
	s.Values[name] = value
}

// Store Get the mongo store
func Store() *gs.CookieStore {
	if _store == nil {
		_store = gs.NewCookieStore(getKey())
	}
	return _store
}

//AddFlash Add a new flash to sessions
func AddFlash(r *http.Request, w http.ResponseWriter, f interface{}, s *gs.Session) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(f)
	if err != nil {
		log.Println(errors.Cause(err))
	}
	s.AddFlash(buf.String(), "Flash")
	err = s.Save(r, w)
	if err != nil {
		log.Println(errors.Cause(err))
	}

}

func GetFlashes(r *http.Request, s *gs.Session) (f []models.Flash) {
	flashmessages := s.Flashes("Flash")
	var flashes []models.Flash
	for _, k := range flashmessages {
		var flash models.Flash
		json.Unmarshal([]byte(k.(string)), &flash)
		flashes = append(flashes, flash)
	}
	return flashes
}

// Signin ...
func (a Admin) Signin(w http.ResponseWriter, r *http.Request) {
	f, err := GetNamed(r, "Flash")
	if err != nil {
		helpers.HandleError(err)
	}
	s, err := GetNamed(r, "current-session")
	if err != nil {
		helpers.HandleError(err)
	}
	if r.Method == "GET" {
		file, err := ioutil.ReadFile("./app/views/admin/signin.html")
		if err != nil {
			helpers.HandleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if s.Values["user-id"] != nil {
			AddFlash(r, w, models.Flash{Type: "Warning", Message: "Already signed in"}, f)
			f.Save(r, w)
			http.Redirect(w, r, "/bd-admin", 302)
			return
		}
		flashes := GetFlashes(r, f)
		t := template.Must(template.New("bd-admin/signin").Funcs(funcMAP).Parse(string(file)))
		err = t.ExecuteTemplate(w, t.Name(), map[string]interface{}{
			"Title":   "Admin Login",
			"Flashes": flashes,
		})
		if err != nil {
			helpers.HandleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	if r.Method == "POST" {
		fmt.Println(r.Form)
		err = r.ParseForm()
		if err != nil {
			helpers.HandleError(err)
		}
		email := r.Form.Get("email")
		fmt.Println(email)
		// sql := fmt.Sprintf("email = %s", email)
		user, err := models.GetUser(fmt.Sprintf("email = '%s'", email))
		if err != nil {
			AddFlash(r, w, models.Flash{Type: "Danger", Message: "No account email/password"}, f)
			f.Save(r, w)
			http.Redirect(w, r, "/bd-admin/signin", 302)
			return
		}

		if user.Email == "" {
			AddFlash(r, w, models.Flash{Type: "Danger", Message: "Incorrect Email"}, f)
			err = f.Save(r, w)
			if err != nil {
				helpers.HandleError(err)
			}
			http.Redirect(w, r, "/bd-admin/signin", 302)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Form.Get("password")))
		if err != nil {
			AddFlash(r, w, models.Flash{Type: "Password", Message: "Incorrect Password"}, f)
			user.LoginAttempts = user.LoginAttempts + 1
			models.UpdateUser(user)
			at := strconv.Itoa(user.LoginAttempts)
			AddFlash(r, w, models.Flash{Type: "Warning", Message: at}, f)
			err = f.Save(r, w)
			if err != nil {
				helpers.HandleError(err)
			}
			http.Redirect(w, r, "/bd-admin/signin", 302)
			return
		}
		if err != nil {
			helpers.HandleError(errors.Trace(err))
		}
		// sum, _ := hasher.NewHasher().SHA1().Write([]byte(string(user.ID)))

		// digest := hash5.DigestString(string(user.ID))

		user.LastLogin = time.Now()
		user.LastIP = helpers.GetIP(r)
		models.UpdateUser(user)
		fmt.Println(user, "<<< where cookie should set")
		s.Values["user-id"] = user.ID
		err = f.Save(r, w)
		if err != nil {
			helpers.HandleError(err)
		}
		s.Save(r, w)
		http.Redirect(w, r, "/bd-admin", 302)
		return
	}
}

//Signout ...
func (a Admin) Signout(w http.ResponseWriter, r *http.Request) {
	s, _ := GetNamed(r, "current-session")
	s.Values["user-id"] = 0
	s.Options.MaxAge = -1
	err := s.Save(r, w)
	if err != nil {
		helpers.HandleError(err)
	}
	s.Save(r, w)
	http.Redirect(w, r, "/", 302)
	return
}

// Signup ... rename function and page and route to register
func (a Admin) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, err := ioutil.ReadFile("./app/views/admin/signup.html")
		if err != nil {
			if err != nil {
				helpers.HandleError(err)
			}
		}
		t := template.Must(template.New("admin").Parse(string(file)))
		err = t.Execute(w, map[string]interface{}{
			"Title": "Admin Signup",
		})
		if err != nil {
			helpers.HandleError(err)
		}
	}
	if r.Method == "POST" {
		s, _ := GetNamed(r, "current-session")
		err := r.ParseForm()
		if err != nil {
			if err != nil {
				helpers.HandleError(err)
			}
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Form.Get("password")), bcrypt.DefaultCost)
		if err != nil {
			helpers.HandleError(err)
		}
		user := models.User{}
		user.FirstName = r.Form.Get("first_name")
		user.LastName = r.Form.Get("last_name")
		user.Email = r.Form.Get("email")
		user.Password = string(hashedPassword)
		err = models.CreateUser(user)
		if err != nil {
			helpers.HandleError(err)
		}
		err = models.CreateUser(user)
		if err != nil {
			helpers.HandleError(err)
		}
		sql := fmt.Sprintf("email = '%s'", user.Email)
		user, err = models.GetUser(sql)
		if err != nil {
			helpers.HandleError(err)
		}
		s.Values["user-id"] = user.ID
		err = s.Save(r, w)
		if err != nil {
			helpers.HandleError(err)
		}
		url := fmt.Sprintf("/bd-admin/user/%d", user.ID)
		http.Redirect(w, r, url, 302)
	}
}
