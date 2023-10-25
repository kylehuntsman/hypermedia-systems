package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/kylehuntsman/hypermedia-systems/contact"
)

type Map map[string]interface{}
type NamedTemplateFiles map[string][]string

type TemplateRenderer struct {
	templateFiles NamedTemplateFiles
	Templates     map[string]*template.Template
}

func NewTemplateRenderer(templateFiles NamedTemplateFiles) *TemplateRenderer {
	return &TemplateRenderer{
		templateFiles: templateFiles,
		Templates:     make(map[string]*template.Template),
	}
}

func (t *TemplateRenderer) Load() error {
	// Load templates
	for name, files := range t.templateFiles {
		tmpl, err := template.New(name).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("failed to parse files for template %s: %w", name, err)
		}
		log.Debugf("Loaded template %s with files %s", tmpl.Name(), tmpl.DefinedTemplates())
		t.Templates[name] = tmpl
	}
	return nil
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, layouts ...string) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return fmt.Errorf("template %s does not exist", name)
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}

// main starts a simple Fiber web server
func main() {
	// Set up contact database
	db := contact.NewContactDB()

	firstContact := contact.NewContact()
	firstContact.FirstName = "Kyle"
	firstContact.LastName = "Huntsman"
	firstContact.Phone = "801-555-1234"
	firstContact.Email = "me@kylehuntsman.com"
	db.AddContact(firstContact)

	secondContact := contact.NewContact()
	secondContact.FirstName = "Kyle"
	secondContact.LastName = "Huntsman"
	secondContact.Phone = "801-555-1234"
	secondContact.Email = "him@kylehuntsman.com"
	db.AddContact(secondContact)

	// Set up web server
	app := fiber.New(fiber.Config{
		Views: NewTemplateRenderer(NamedTemplateFiles{
			"contacts.html":     []string{"./templates/contacts.html", "./templates/base.html"},
			"contact-new.html":  []string{"./templates/contact-new.html", "./templates/base.html"},
			"contact-view.html": []string{"./templates/contact-view.html", "./templates/base.html"},
			"contact-edit.html": []string{"./templates/contact-edit.html", "./templates/base.html"},
		}),
	})

	app.Static("/", "./public")

	// Log all requests
	log.SetLevel(log.LevelDebug)
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts", fiber.StatusFound)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		searchQuery := c.Query("q")

		var contacts []*contact.Contact
		if searchQuery == "" {
			contacts = db.GetAllContacts()
		} else {
			contacts = db.SearchContacts(searchQuery)
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Render("contacts.html", Map{
			"Contacts":    contacts,
			"SearchQuery": searchQuery,
		})
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		return c.Render("contact-new.html", contact.NewContact())
	})
	app.Post("/contacts/new", func(c *fiber.Ctx) error {
		contact := contact.NewContact()
		contact.FirstName = utils.CopyString(c.FormValue("first_name"))
		contact.LastName = utils.CopyString(c.FormValue("last_name"))
		contact.Phone = utils.CopyString(c.FormValue("phone"))
		contact.Email = utils.CopyString(c.FormValue("email"))

		ok := db.AddContact(contact)
		if !ok {
			return c.Render("contact-new.html", contact)
		}

		// Created new contact!
		return c.Redirect("/contacts", fiber.StatusFound)
	})

	app.Get("/contacts/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Errorf("Failed to parse id %s url parameter: %w", idParam, err)
		}
		log.Debugw("Parsed id", "id", id)

		contact, ok := db.GetContactById(id)
		if !ok {
			log.Errorf("Failed to get contact by id: id %s not found", idParam)
		}

		return c.Render("contact-view.html", contact)
	})

	app.Get("/contacts/:id/edit", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Errorf("Failed to parse id %s url parameter: %w", idParam, err)
		}
		log.Debugw("Parsed id", "id", id)

		contact, ok := db.GetContactById(id)
		if !ok {
			log.Errorf("Failed to get contact by id: id %s not found", idParam)
		}

		return c.Render("contact-edit.html", contact)
	})
	app.Post("/contacts/:id/edit", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Errorf("Failed to parse id %s url parameter: %w", idParam, err)
		}
		log.Debugw("Parsed id", "id", id)

		contact := contact.NewContact()
		contact.Id = id
		contact.FirstName = utils.CopyString(c.FormValue("first_name"))
		contact.LastName = utils.CopyString(c.FormValue("last_name"))
		contact.Phone = utils.CopyString(c.FormValue("phone"))
		contact.Email = utils.CopyString(c.FormValue("email"))

		ok := db.UpdateContact(contact)
		if !ok {
			return c.Render("contact-edit.html", contact)
		}

		// Created new contact!
		return c.Redirect(fmt.Sprintf("/contacts/%s", idParam), fiber.StatusFound)
	})

	app.Delete("/contacts/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Errorf("Failed to parse id %s url parameter: %w", idParam, err)
		}
		log.Debugw("Parsed id", "id", id)

		ok := db.DeleteContactById(id)
		if !ok {
			log.Errorf("Failed to delete contact %s: id not found", idParam)
		}

		return c.Redirect("/contacts", fiber.StatusSeeOther)
	})

	app.Get("/contacts/:id/edit/email", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Errorf("Failed to parse id %s url parameter: %w", idParam, err)
		}
		log.Debugw("Parsed id", "id", id)

		contact, ok := db.GetContactById(id)
		if !ok {
			log.Errorf("Failed to get contact by id: id %s not found", idParam)
		}

		copy := contact.Copy()
		email := utils.CopyString(c.FormValue("email"))
		copy.Email = email
		if db.HasEmail(copy.Id, copy.Email) {
			copy.Errors["Email"] = "Email must be unique"
		}

		return c.Render("contact-edit.html", copy)
	})

	app.Listen(":3000")
}
