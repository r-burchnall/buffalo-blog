package actions

func (as *ActionSuite) Test_HomeHandler() {
	res := as.HTML("/health").Get()
	as.Equal(200, res.Code)
}
