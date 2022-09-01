package grifts

import (
	"github.com/gobuffalo/grift/grift"
	"log"
	"social_api/models"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		defaultUser := &models.User{
			Email:    "test@test.com",
			Password: "test",
		}

		err := models.DB.Create(defaultUser)
		if err != nil {
			log.Fatalf("failed to create user in database, %v", err)
		}

		return nil
	})

})
