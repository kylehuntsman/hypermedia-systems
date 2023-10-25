package contact

import "github.com/google/uuid"

type Contact struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Errors    map[string]string
}

func NewContact() *Contact {
	return &Contact{
		Id:     uuid.New(),
		Errors: map[string]string{},
	}
}

func EmptyContact() *Contact {
	return &Contact{
		Id:     uuid.Nil,
		Errors: map[string]string{},
	}
}

func (c *Contact) Copy() *Contact {
	return &Contact{
		Id:        c.Id,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Phone:     c.Phone,
		Email:     c.Email,
		Errors:    map[string]string{},
	}
}

type ContactDB struct {
	contacts []*Contact
}

func NewContactDB() *ContactDB {
	return &ContactDB{[]*Contact{}}
}

func (db *ContactDB) AddContact(contact *Contact) bool {
	if db.HasEmail(contact.Id, contact.Email) {
		contact.Errors["Email"] = "Email must be unique"
		return false
	}

	db.contacts = append(db.contacts, contact)
	return true
}

func (db *ContactDB) GetContactById(id uuid.UUID) (*Contact, bool) {
	for _, c := range db.contacts {
		if c.Id == id {
			return c, true
		}
	}
	return EmptyContact(), false
}

func (db *ContactDB) HasEmail(id uuid.UUID, email string) bool {
	for _, c := range db.contacts {
		if c.Id != id && c.Email == email {
			return true
		}
	}
	return false
}

func (db *ContactDB) GetAllContacts() []*Contact {
	return db.contacts
}

func (db *ContactDB) SearchContacts(search string) []*Contact {
	return db.contacts
}

func (db *ContactDB) DeleteContactById(id uuid.UUID) bool {
	for i, c := range db.contacts {
		if c.Id == id {
			db.contacts = append(db.contacts[:i], db.contacts[i+1:]...)
			return true
		}
	}

	return false
}

func (db *ContactDB) UpdateContact(contact *Contact) bool {
	if db.HasEmail(contact.Id, contact.Email) {
		contact.Errors["Email"] = "Email must be unique"
		return false
	}

	for i, c := range db.contacts {
		if c.Id == contact.Id {
			db.contacts[i] = contact
			return true
		}
	}

	return false
}
