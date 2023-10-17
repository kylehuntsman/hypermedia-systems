package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/kylehuntsman/hypermedia-systems/contact"
	"github.com/kylehuntsman/hypermedia-systems/templates" // If you're having a broken import error, run make prepare to generate the template files
)

// main starts a simple Fiber web server
func main() {
	db := contact.NewContactDB()

	app := fiber.New()

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
		contactsView := templates.Contacts(contacts, searchQuery)
		return templates.Index(contactsView).Render(context.Background(), c)
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		newContactView := templates.NewContact(contact.NewContact())
		return templates.Index(newContactView).Render(context.Background(), c)
	})
	app.Post("/contacts/new", func(c *fiber.Ctx) error {
		contact := contact.NewContact()
		contact.FirstName = utils.CopyString(c.FormValue("first_name"))
		contact.LastName = utils.CopyString(c.FormValue("last_name"))
		contact.Phone = utils.CopyString(c.FormValue("phone"))
		contact.Email = utils.CopyString(c.FormValue("email"))

		ok := db.AddContact(contact)

		if !ok {
			c.Set("Content-Type", "text/html; charset=utf-8")
			newContactView := templates.NewContact(contact)
			return templates.Index(newContactView).Render(context.Background(), c)
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

		c.Set("Content-Type", "text/html; charset=utf-8")
		viewContactView := templates.ViewContact(contact)
		return templates.Index(viewContactView).Render(context.Background(), c)
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

		c.Set("Content-Type", "text/html; charset=utf-8")
		editContactView := templates.EditContact(contact)
		return templates.Index(editContactView).Render(context.Background(), c)
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
		contact.FirstName = c.FormValue("first_name")
		contact.LastName = c.FormValue("last_name")
		contact.Phone = c.FormValue("phone")
		contact.Email = c.FormValue("email")

		ok := db.UpdateContact(contact)
		if !ok {
			c.Set("Content-Type", "text/html; charset=utf-8")
			editContactView := templates.EditContact(contact)
			return templates.Index(editContactView).Render(context.Background(), c)
		}

		// Created new contact!
		return c.Redirect(fmt.Sprintf("/contacts/%s", idParam), fiber.StatusFound)
	})

	app.Post("/contacts/:id/delete", func(c *fiber.Ctx) error {
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

		return c.Redirect("/contacts", fiber.StatusFound)
	})

	app.Listen(":3000")
}
