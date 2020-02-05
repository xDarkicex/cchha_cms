package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"regexp"
	"time"

	e "github.com/scorredoira/email"

	"github.com/alecthomas/template"
	"github.com/juju/errors"
	"github.com/xDarkicex/cchha_server_new/app/models"
	"github.com/xDarkicex/cchha_server_new/helpers"
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type Mailer Controllers

func (mailer Mailer) Career(res http.ResponseWriter, req *http.Request) {
	f, _ := GetNamed(req, "Flash")
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		helpers.HandleError(err)
	}

	file, handle, err := req.FormFile("file_upload")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	defer file.Close()
	fi, err := os.OpenFile("./test/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer file.Close()
	_, err = io.Copy(fi, file)
	if err != nil {
		helpers.HandleError(err)
	}
	bytefile, err := ioutil.ReadAll(file)
	if err != nil {
		helpers.HandleError(err)
	}
	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "application/pdf":
		err := mailer.SaveFile(handle.Filename, bytefile)
		if err != nil {
			AddFlash(req, res, models.Flash{Type: "Warn", Message: "Error saving file"}, f)
			helpers.HandleError(errors.Trace(err))
		}
	default:
		AddFlash(req, res, models.Flash{Type: "Warn", Message: "Unsupported File Type"}, f)
		http.Redirect(res, req, "/home-health/careers.html", 302)
		return
	}

	name := (req.FormValue("contact_name"))
	sender := (req.FormValue("contact_email"))
	add := (req.FormValue("contact_address"))
	phone := (req.FormValue("contact_phone"))
	// body := (req.FormValue("contact_body"))
	google_chaptcha := (req.FormValue("g-recaptcha-response"))

	google_struct := validate_google_rechaptcha(google_chaptcha)
	if !google_struct.Success {
		AddFlash(req, res, models.Flash{Type: "Warn", Message: "Failed reChaptcha"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(err)
		}
		http.Redirect(res, req, "/home-health/careers.html", 302)
		return
	}
	if !validate_email(sender) {
		AddFlash(req, res, models.Flash{Type: "Warn", Message: "Invalid Email"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(err)
		}
		http.Redirect(res, req, "/home-health/careers.html", 302)
		return
	}

	t, err := template.ParseFiles("./app/views/email/job.html")
	if err != nil {
		helpers.HandleError(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, map[string]interface{}{
		"Phone":   phone,
		"Address": add,
		"Sender":  name,
		"Email":   sender,
		"Content": (req.FormValue("contact_body")),
	}); err != nil {
		helpers.HandleError(err)
	}
	subject := fmt.Sprintf("Resume submission cchha")
	msg := buf.String()
	m := e.NewHTMLMessage(subject, msg)
	m.From = mail.Address{Name: "Jobs", Address: "admin@cchha.com"}
	// testing jobs@cchha.com
	m.To = []string{"xDarkicex@gmail.com"}
	err = m.Attach("/tmp/" + handle.Filename)
	if err != nil {
		helpers.HandleError(err)
	}
	// HomeHealth2017
	auth := smtp.PlainAuth("", "admin@cchha.com", "Vh2@cchha#G0!", "smtp.gmail.com")
	SMTP := "smtp.gmail.com:587"
	if err := e.Send(SMTP, auth, m); err != nil {
		helpers.HandleError(errors.Trace(err))
	}

	AddFlash(req, res, models.Flash{Type: "Success", Message: "Email Sent"}, f)
	err = f.Save(req, res)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	mailer.DeleteFile(handle.Filename)
	http.Redirect(res, req, "/home-health/contact-careers", 302)
}

func (mailer Mailer) Contact(res http.ResponseWriter, req *http.Request) {
	fmt.Println(helpers.Red("Email"), "values hardcoded")
	f, _ := GetNamed(req, "Flash")
	senderEmail := (req.FormValue("contact_email"))
	add := (req.FormValue("contact_address"))
	phone := (req.FormValue("contact_phone"))
	Sender := (req.FormValue("contact_name"))
	google_chaptcha := (req.FormValue("g-recaptcha-response"))

	google_struct := validate_google_rechaptcha(google_chaptcha)
	if !google_struct.Success {
		AddFlash(req, res, models.Flash{Type: "Failure", Message: "Failed reChaptcha"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(errors.Trace(err))
		}
		http.Redirect(res, req, "/contact.html", 302)
		return
	}
	if !validate_email(senderEmail) {
		AddFlash(req, res, models.Flash{Type: "Failure", Message: "Not Valid Email"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(errors.Trace(err))
		}
		http.Redirect(res, req, "/contact.html", 302)
		return
	}
	// Start of html email sending

	subject := fmt.Sprintf("New Contact Request from %s", Sender)
	t, err := template.ParseFiles("./app/views/email/contact.html")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, map[string]interface{}{
		"Phone":   phone,
		"Address": add,
		"Sender":  Sender,
		"Email":   senderEmail,
		"Content": (req.FormValue("contact_body")),
	}); err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	Body := buf.String()
	msg := Body
	m := e.NewHTMLMessage(subject, msg)
	m.Subject = subject
	m.BodyContentType = "text/html"
	m.From = mail.Address{Name: "Compassionate Care Home Health", Address: "admin@cchha.com"}
	// testing jobs@cchha.com
	m.To = []string{"xDarkicex@gmail.com"}
	// HomeHealth2017
	auth := smtp.PlainAuth("", "admin@cchha.com", "Vh2@cchha#G0!", "smtp.gmail.com")
	SMTP := "smtp.gmail.com:587"
	if err := e.Send(SMTP, auth, m); err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	AddFlash(req, res, models.Flash{Type: "Success", Message: "Email sent successfully"}, f)
	err = f.Save(req, res)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
}

func (mailer Mailer) Reviews(res http.ResponseWriter, req *http.Request) {

	f, _ := GetNamed(req, "Flash")
	Recipient := (req.FormValue("contact_name"))
	Sender := (req.FormValue("contact_sender"))
	Phone := (req.FormValue("contact_sender_phone"))
	JobTitle := (req.FormValue("contact_sender_title"))
	Email := (req.FormValue("contact_email"))
	Body := (req.FormValue("contact_body"))
	GoogleChaptcha := (req.FormValue("g-recaptcha-response"))
	GoogleStruct := validate_google_rechaptcha(GoogleChaptcha)

	if !GoogleStruct.Success {
		AddFlash(req, res, models.Flash{Type: "Failure", Message: "Failed reChaptcha"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(errors.Trace(err))
		}
		http.Redirect(res, req, "/bd-admin", 302)
		return
	}
	if !validate_email(Email) {
		AddFlash(req, res, models.Flash{Type: "Failure", Message: "Must enter valid email"}, f)
		err := f.Save(req, res)
		if err != nil {
			helpers.HandleError(errors.Trace(err))
		}
		http.Redirect(res, req, "/bd-admin", 302)
		return
	}
	t, err := template.ParseFiles("./app/views/email/review.html")
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, map[string]interface{}{
		"Sender":    Sender,
		"Recipient": Recipient,
		"Email":     Email,
		"Phone":     Phone,
		"Title":     JobTitle,
		"Content":   Body,
	}); err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	subject := fmt.Sprintf("Thank you %s", Recipient)
	msg := buf.String()
	m := e.NewHTMLMessage(subject, msg)

	m.Subject = subject

	m.From = mail.Address{Name: "Compassionate Care Home Health", Address: "admin@cchha.com"}
	// testing jobs@cchha.com
	m.To = []string{"xDarkicex@gmail.com"}
	// HomeHealth2017

	auth := smtp.PlainAuth("", "admin@cchha.com", "Vh2@cchha#G0!", "smtp.gmail.com")
	SMTP := "smtp.gmail.com:587"
	if err := e.Send(SMTP, auth, m); err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	AddFlash(req, res, models.Flash{Type: "Success", Message: "Email sent successfully"}, f)
	err = f.Save(req, res)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	http.Redirect(res, req, "/bd-admin", 302)
}

func validate_email(email string) bool {
	regex, err := regexp.Compile(`\S+@\S+`)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	if !regex.MatchString(email) {
		return false
	}
	return true
}

func validate_google_rechaptcha(chaptcha string) (r RecaptchaResponse) {
	google_check := url.Values{
		"secret":   {"6Ldyv7UUAAAAAIv9q4YMxZuBPqdmM1k6m4rI4HN1"},
		"response": {chaptcha},
		"remoteip": {"127.0.0.1"},
	}
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", google_check)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	defer resp.Body.Close()
	google_body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Read error: could not read body: %s", err)
	}
	err = json.Unmarshal(google_body, &r)
	if err != nil {
		log.Printf("Read error: got invalid JSON: %s", err)
	}
	// fmt.Println(r)
	return r
}

func GenerateCookie(status string, success bool) *http.Cookie {
	type data struct {
		Status  string
		Success bool
	}
	cookie_value := data{
		Status:  status,
		Success: success,
	}
	d, err := json.Marshal(cookie_value)
	if err != nil {
		helpers.HandleError(errors.Trace(err))
	}
	fmt.Println(string(d))
	cookie := &http.Cookie{
		Name:  "message",
		Value: base64.StdEncoding.EncodeToString(d),
		// Path:    "cchha.com/contact.html",
		// Domain:  "cchha.com",
		Expires: time.Now().Add(time.Minute * 5),

		Secure:   false,
		HttpOnly: false,
	}
	// RawExpires: "0",
	// MaxAge:     0,
	return cookie
}

// Bad func used for internal purposes
func (mailer Mailer) SaveFile(fileName string, file []byte) (err error) {
	err = ioutil.WriteFile("/tmp/"+fileName, file, 0666)
	if err != nil {
		helpers.HandleError(err)
		// http.Error(, errors.Details(err), http.StatusInternalServerError)
	}
	// db.Save(file)
	return nil
}

func (mailer Mailer) DeleteFile(fileName string) {
	//	os.Remove("/tmp/" + fileName)
}
