package main

import (
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"path/filepath"
)

type PageData struct {
	Title     string
	BodyClass string
	Active    string
	Success   string
	Error     string
	Project   string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	layout := "templates/layout.html"
	page := filepath.Join("templates", tmpl)

	t, err := template.ParseFiles(layout, page)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", PageData{Title: "Welcome - Ajay Pandey", BodyClass: "index", Active: "home"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about.html", PageData{Title: "About - Ajay Pandey", BodyClass: "about-page", Active: "about"})
	})

	http.HandleFunc("/resume", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "resume.html", PageData{Title: "Resume - Ajay Pandey", BodyClass: "resume-page", Active: "resume"})
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "services.html", PageData{Title: "Services - Ajay Pandey", BodyClass: "services-page", Active: "services"})
	})

	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "portfolio.html", PageData{Title: "Portfolio - Ajay Pandey", BodyClass: "portfolio-page", Active: "portfolio"})
	})

	http.HandleFunc("/portfolio-details", func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project") // get the project key, e.g., "hrms"
		data := PageData{
			Title:     "Portfolio Details - Ajay Pandey",
			BodyClass: "portfolio-details-page",
			Active:    "portfolio",
			Project:   project,
		}
		renderTemplate(w, "portfolio-details.html", data)
	})

	// ✅ Contact Page GET + POST
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			email := r.FormValue("email")
			subject := r.FormValue("subject")
			message := r.FormValue("message")

			// Gmail SMTP settings
			from := "poojaprajapati269@gmail.com" // apna Gmail daalo
			password := "yhklkoqcipadycss"        // Gmail App Password generate karo (normal password nahi chalega)

			to := []string{"apandey3212@gmail.com"}
			smtpHost := "smtp.gmail.com"
			smtpPort := "587"

			// Email body
			body := "From: " + name + " <" + email + ">\n" +
				"To: " + to[0] + "\n" +
				"Subject: " + subject + "\n\n" +
				message

			auth := smtp.PlainAuth("", from, password, smtpHost)
			err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(body))

			if err != nil {
				log.Println("Error sending email:", err)
				renderTemplate(w, "contact.html", PageData{
					Title:     "Contact - Ajay Pandey",
					BodyClass: "contact-page",
					Active:    "contact",
					Error:     "Message not sent. Please try again.",
				})
				return
			}

			renderTemplate(w, "contact.html", PageData{
				Title:     "Contact - Ajay Pandey",
				BodyClass: "contact-page",
				Active:    "contact",
				Success:   "Your message has been sent successfully!",
			})
			return
		}

		// GET request
		renderTemplate(w, "contact.html", PageData{Title: "Contact - Ajay Pandey", BodyClass: "contact-page", Active: "contact"})
	})

	log.Println("✅ Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
