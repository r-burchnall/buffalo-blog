package models

func (ms *ModelSuite) Test_Post() {

	ms.Run("should validate a post", func() {
		newPost := Post{
			Content: "this is content",
		}
		verrs, err := newPost.Validate(ms.DB)
		assertNoValidationErrors(ms.T(), verrs)
		assertError(ms.T(), err, nil)
	})

	ms.Run("should fail if missing content", func() {
		newPost := Post{}
		verrs, err := newPost.Validate(ms.DB)
		ms.True(verrs.HasAny())

		contentError := verrs.Get("content")
		ms.Contains(contentError, "Content field is required")

		assertError(ms.T(), err, nil)
	})
}
