package models

func (ms *ModelSuite) Test_Feed() {
	type feedTest struct {
		Name          string
		model         Feed
		expectSuccess bool
	}

	tests := []feedTest{
		{
			"valid feed",
			Feed{
				Name: "A brilliant feed",
				Type: UserProfileFeed,
			},
			true,
		},
		{
			"missing feed name",
			Feed{
				Type: UserProfileFeed,
			},
			false,
		},
	}

	for _, test := range tests {
		ms.Run(test.Name, func() {
			verrs, err := test.model.Validate(ms.DB)
			if err != nil {
				ms.Failf("failure to call validate, %v", err.Error())
			}

			// Flip semantic value for test suite
			ms.Equalf(verrs.HasAny(), !test.expectSuccess, "validation errors %v", verrs)
		})
	}
}
