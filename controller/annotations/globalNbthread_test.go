package annotations

import (
	"github.com/haproxytech/kubernetes-ingress/controller/store"
)

func (suite *AnnotationSuite) TestNbthreadUpdate() {
	test := store.StringW{Value: "1"}
	a := NewGlobalNbthread("", suite.client)
	if suite.NoError(a.Parse(test, true)) {
		suite.NoError(a.Update())
		result, _ := suite.client.GlobalWriteConfig("global", "nbthread")
		suite.Equal("nbthread 1", result)
	}
}

func (suite *AnnotationSuite) TestNbthreadFail() {
	test := store.StringW{Value: "garbage"}
	a := NewGlobalNbthread("", suite.client)
	err := a.Parse(test, true)
	suite.T().Log(err)
	suite.Error(err)
}
