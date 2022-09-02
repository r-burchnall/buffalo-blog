package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"social_api/models"
)

func (as *ActionSuite) Test_FeedsResource_List() {
	as.LoadFixture("lots of posts")
	response := as.JSON("/feeds").Get()
	as.True(response.Code == http.StatusOK)

	var feeds models.Feeds
	err := json.NewDecoder(response.Body).Decode(&feeds)
	if err != nil {
		as.Fail("failed to parse body")
	}

	if len(feeds) != 1 {
		as.Fail("scenario should return 2 feeds")
	}
}

func (as *ActionSuite) Test_FeedsResource_Show() {
	as.LoadFixture("lots of posts")

	firstFeed, _ := getFirstFeed(as)

	response := as.JSON(fmt.Sprintf("/feeds/%s", firstFeed.ID)).Get()
	as.True(response.Code == http.StatusOK)

	firstFeedJsonBuffer := bytes.Buffer{}
	err := json.NewEncoder(&firstFeedJsonBuffer).Encode(&firstFeed)
	if err != nil {
		as.Fail("fail to encode json of first feed")
	}

	as.JSONEq(response.Body.String(), firstFeedJsonBuffer.String())
}

func (as *ActionSuite) Test_FeedsResource_Create() {
	as.Run("create a valid feed", func() {
		newFeed := models.Feed{
			Name: "example feed",
			Type: models.UserProfileFeed,
		}
		response := as.JSON("/feeds").Post(&newFeed)
		as.Equal(response.Code, http.StatusCreated, "Unexpected response %s", response.Body.String())
	})

	as.Run("return validation errors for an invalid feed", func() {
		newFeed := models.Feed{}
		response := as.JSON("/feeds").Post(&newFeed)
		as.Equal(response.Code, http.StatusUnprocessableEntity)

		as.Contains(response.Body.String(), "name")
	})

}

func (as *ActionSuite) Test_FeedsResource_Update() {
	as.Run("should update existing feed", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")

		firstFeed, _ := getFirstFeed(as)
		firstFeed.Name = "EDIT: updated content"

		response := as.JSON(fmt.Sprintf("/feeds/%s", firstFeed.ID)).Put(&firstFeed)
		as.Equal(response.Code, http.StatusOK)

		refreshedFeed := models.Feed{}
		err := as.DB.Q().Where("id = ?", firstFeed.ID).First(&refreshedFeed)
		if err != nil {
			as.Fail("fail to get feed after update")
		}

		var responseFeed models.Feed
		err = json.NewDecoder(response.Body).Decode(&responseFeed)
		if err != nil {
			as.Fail("failed to marshal response from update")
		}
		as.Equal(responseFeed.Name, refreshedFeed.Name)
	})

	as.Run("should fail to update none existent feed", func() {
		as.SetupTest()
		as.LoadFixture("lots of posts")

		firstFeed, _ := getFirstFeed(as)
		originalName := firstFeed.Name
		firstFeed.Name = "EDIT: updated name"

		response := as.JSON(fmt.Sprintf("/feeds/%s", "bad-id")).Put(&firstFeed)
		as.Equal(response.Code, http.StatusNotFound)

		unchangedFeed := models.Feed{}
		err := as.DB.Q().Where("id = ?", firstFeed.ID).First(&unchangedFeed)
		if err != nil {
			as.Fail("fail to get feed after update")
		}

		as.Equal(originalName, unchangedFeed.Name)
	})
}

func (as *ActionSuite) Test_FeedsResource_Destroy() {
	as.SetupTest()

	as.Run("delete an existing feed", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")
		firstFeed, _ := getFirstFeed(as)
		total, err := as.DB.Count(&firstFeed)
		if err != nil {
			as.Fail("could not get total feeds")
		}
		as.Equal(total, 1)

		response := as.JSON(fmt.Sprintf("/feeds/%s", firstFeed.ID)).Delete()
		as.Equal(response.Code, http.StatusOK)

		total, err = as.DB.Count(&firstFeed)
		if err != nil {
			as.Fail("could not get total feeds after delete")
		}
		as.Equal(total, 0)
	})

	as.Run("not found response", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")
		response := as.JSON(fmt.Sprintf("/feeds/%s", "bad-id")).Delete()
		as.Equal(response.Code, http.StatusNotFound)
	})
}

// Gets the first feed in the database, will Fail test if it cannot find a feed
func getFirstFeed(as *ActionSuite) (firstFeed models.Feed, err error) {
	err = models.DB.First(&firstFeed)
	if err != nil {
		as.Failf("failed to load first feed from db: %s", err.Error())
	}
	return
}
