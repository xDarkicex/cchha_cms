package server

import (
	"net/http"

	"github.com/xDarkicex/cchha_server_new/helpers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/xDarkicex/cchha_server_new/app/controllers"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.StripSlashes)
	router.Use(helpers.Elephant)

	// Routes //

	// splash
	application := controllers.Application{}
	router.Get("/", application.Index)
	router.NotFound(application.CustomNotFound)
	// auth middleware controller binding
	auth := controllers.Auth{}

	// webhook

	// router.Route("/api", func(r chi.Router) {
	// 	r.Get("/", DialSocket)
	// })
	// admin panel things
	admin := controllers.Admin{}
	router.Get("/bd-admin/signin", admin.Signin)
	router.Post("/bd-admin/signin", admin.Signin)
	router.Get("/bd-admin/signout", auth.Authtenicate(admin.Signout))
	router.Get("/bd-admin/signup", auth.Authtenicate(admin.Signup))
	router.Post("/bd-admin/signup", admin.Signup)
	router.Get("/bd-admin/user/{userID}/edit", auth.Authtenicate(admin.Edit))
	router.Post("/bd-admin/user/{userID}/edit", admin.Edit)
	router.Get("/bd-admin/user/{userID:\\d{1,200}$}", admin.Show)
	router.Get("/bd-admin", auth.Authtenicate(admin.Index))
	router.Post("/bd-admin/review/create", admin.Create)
	// router.Get("/bd-admin/user/forgot", admin.Forgot)
	// router.Get("/bd-admin/user/reset", admin.Reset)
	// router.Get("/bd-admin/user/reset/token={value:^.{24}$}", admin.Reset)
	// router.Post("/bd-admin/user/forgot", admin.Forgot)

	// review

	// home health
	homehealth := controllers.HomeHealth{}
	router.Get("/home-health", homehealth.Index)
	router.Get("/home-health/careers", homehealth.Careers)
	router.Get("/home-health/careers.html", homehealth.Careers)
	router.Get("/home-health/services", homehealth.Services)
	router.Get("/home-health/services.html", homehealth.Services)
	router.Get("/home-health/eligibility", homehealth.Eligibility)
	router.Get("/home-health/eligibility.html", homehealth.Eligibility)
	router.Get("/home-health/resources", homehealth.Resources)
	router.Get("/home-health/resources.html", homehealth.Resources)
	router.Get("/home-health/community", homehealth.Community)
	router.Get("/home-health/community.html", homehealth.Community)
	router.Get("/home-health/about", homehealth.About)
	router.Get("/home-health/about.html", homehealth.About)
	router.Get("/home-health/locations", homehealth.Locations)
	router.Get("/home-health/locations.html", homehealth.Locations)
	router.Get("/home-health/contact", homehealth.Contact)
	router.Get("/home-health/contact.html", homehealth.Contact)

	// reviews
	reviews := controllers.Reviews{}
	router.Get("/home-health/reviews", reviews.Index)
	router.Get("/home-health/reviews?", reviews.Index)
	router.Get("/home-health/reviews/json", reviews.Json)
	router.Get("/home-health/reviews.html", reviews.Index)
	router.Post("/home-health/reviews", reviews.Create)
	router.Get("/home-health/review/{reviewID}", reviews.Show)
	router.Get("/home-health/review/{reviewID}/details", reviews.Details)
	router.Get("/home-health/reviews/{sort}", reviews.Sort)

	// hospice
	hospice := controllers.Hospice{}
	router.Get("/hospice", hospice.Index)
	router.Get("/hospice/careers", hospice.Careers)
	router.Get("/hospice/careers.html", hospice.Careers)
	router.Get("/hospice/services", hospice.Services)
	router.Get("/hospice/services.html", hospice.Services)
	router.Get("/hospice/eligibility", hospice.Eligibility)
	router.Get("/hospice/eligibility.html", hospice.Eligibility)
	router.Get("/hospice/resources", hospice.Resources)
	router.Get("/hospice/resources.html", hospice.Resources)
	router.Get("/hospice/community", hospice.Community)
	router.Get("/hospice/community.html", hospice.Community)
	router.Get("/hospice/about", hospice.About)
	router.Get("/hospice/about.html", hospice.About)
	router.Get("/hospice/locations", hospice.Locations)
	router.Get("/hospice/locations.html", hospice.Locations)
	router.Get("/hospice/contact", hospice.Contact)
	router.Get("/hospice/contact.html", hospice.Contact)

	// Post request handing google email settings
	mail := controllers.Mailer{}
	router.Post("/contact-careers", mail.Career)
	router.Post("/contact", mail.Contact)
	router.Post("/bd-admin/review/response", mail.Reviews)

	// // Api calls
	router.Post("/bd-admin/review/reject", admin.Reject)
	router.Post("/bd-admin/review/approve", admin.Approve)
	router.Post("/bd-admin/review/delete", admin.DeleteReview)
	router.Post("/bd-admin/review/feature", reviews.Feature)

	// static file server
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./public/")))
	router.Handle("/static/*", fileServer)

	return router
}
