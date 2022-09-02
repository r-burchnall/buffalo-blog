package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"social_api/models"
)

func (as *ActionSuite) Test_PostsResource_List() {
	as.LoadFixture("lots of posts")
	response := as.JSON("/posts").Get()
	as.True(response.Code == http.StatusOK)

	var posts models.Posts
	err := json.NewDecoder(response.Body).Decode(&posts)
	if err != nil {
		as.Fail("failed to parse body")
	}

	if len(posts) != 2 {
		as.Fail("scenario should return 2 posts")
	}
}

func (as *ActionSuite) Test_PostsResource_Show() {
	as.LoadFixture("lots of posts")

	firstPost, _ := getFirstPost(as)

	response := as.JSON(fmt.Sprintf("/posts/%s", firstPost.ID)).Get()
	as.True(response.Code == http.StatusOK)

	firstPostJsonBuffer := bytes.Buffer{}
	err := json.NewEncoder(&firstPostJsonBuffer).Encode(&firstPost)
	if err != nil {
		as.Fail("fail to encode json of first post")
	}

	as.JSONEq(response.Body.String(), firstPostJsonBuffer.String())
}

func (as *ActionSuite) Test_PostsResource_Create() {
	as.Run("create a valid post", func() {
		newPost := models.Post{
			Content: "this is some content",
		}
		response := as.JSON("/posts").Post(&newPost)
		as.Equal(response.Code, http.StatusCreated)
	})

	as.Run("return validation errors for an invalid post", func() {
		newPost := models.Post{}
		response := as.JSON("/posts").Post(&newPost)
		as.Equal(response.Code, http.StatusUnprocessableEntity)

		as.Contains(response.Body.String(), "content")
	})

}

func (as *ActionSuite) Test_PostsResource_Update() {
	as.Run("should update existing post", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")

		firstPost, _ := getFirstPost(as)
		firstPost.Content = "EDIT: updated content"

		response := as.JSON(fmt.Sprintf("/posts/%s", firstPost.ID)).Put(&firstPost)
		as.Equal(response.Code, http.StatusOK)

		refreshedPost := models.Post{}
		err := as.DB.Q().Where("id = ?", firstPost.ID).First(&refreshedPost)
		if err != nil {
			as.Fail("fail to get post after update")
		}

		var responsePost models.Post
		err = json.NewDecoder(response.Body).Decode(&responsePost)
		if err != nil {
			as.Fail("failed to marshal response from update")
		}
		as.Equal(responsePost.Content, refreshedPost.Content)
	})

	as.Run("should fail to update none existent post", func() {
		as.SetupTest()
		as.LoadFixture("lots of posts")

		firstPost, _ := getFirstPost(as)
		originalContent := firstPost.Content
		firstPost.Content = "EDIT: updated content"

		response := as.JSON(fmt.Sprintf("/posts/%s", "bad-id")).Put(&firstPost)
		as.Equal(response.Code, http.StatusNotFound)

		unchangedPost := models.Post{}
		err := as.DB.Q().Where("id = ?", firstPost.ID).First(&unchangedPost)
		if err != nil {
			as.Fail("fail to get post after update")
		}

		as.Equal(originalContent, unchangedPost.Content)
	})
}

func (as *ActionSuite) Test_PostsResource_Destroy() {
	as.SetupTest()

	as.Run("delete an existing post", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")
		firstPost, _ := getFirstPost(as)
		total, err := as.DB.Count(&firstPost)
		if err != nil {
			as.Fail("could not get total posts")
		}
		as.Equal(total, 2)

		response := as.JSON(fmt.Sprintf("/posts/%s", firstPost.ID)).Delete()
		as.Equal(response.Code, http.StatusOK)

		total, err = as.DB.Count(&firstPost)
		if err != nil {
			as.Fail("could not get total posts after delete")
		}
		as.Equal(total, 1)
	})

	as.Run("not found response", func() {
		as.SetupTest()

		as.LoadFixture("lots of posts")
		response := as.JSON(fmt.Sprintf("/posts/%s", "bad-id")).Delete()
		as.Equal(response.Code, http.StatusNotFound)
	})
}

// Gets the first post in the database, will Fail test if it cannot find a post
func getFirstPost(as *ActionSuite) (firstPost models.Post, err error) {
	err = models.DB.First(&firstPost)
	if err != nil {
		as.Failf("failed to load first post from db: %s", err.Error())
	}
	return
}
