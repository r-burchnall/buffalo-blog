package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Feed is used by pop to map your feeds database table to your go code.
type Feed struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Type        FeedType     `json:"type" db:"type"`
	Name        string       `json:"name" db:"name"`
	Description nulls.String `json:"description" db:"description"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

type FeedType int

const (
	Unknown FeedType = iota
	UserProfileFeed
	WritersBlock
	Community
	Game
	Gallery
)

// String is not required by pop and may be deleted
func (f Feed) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Feeds is not required by pop and may be deleted
type Feeds []Feed

// String is not required by pop and may be deleted
func (f Feeds) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Feed) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(f.Type), Name: "Type"},
		&validators.StringIsPresent{Field: f.Name, Name: "Name"},
		&validators.FuncValidator{Fn: func() bool {
			switch f.Type {
			case UserProfileFeed:
				return true
			case WritersBlock:
				return true
			case Community:
				return true
			case Game:
				return true
			case Gallery:
				return true
			default:
				return false
			}
		}, Field: "FeedType", Name: "Type Value", Message: "Unexpected value for feed type"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Feed) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Feed) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
