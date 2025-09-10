package v1

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service/v1/login"
)

type V1Group interface {
	login.LoginItf
}

type v1 struct {
	login.LoginItf
}

func GetV1Group(dao dao.RepositoryItf) V1Group {
	return &v1{
		login.GetLoginItf(dao),
	}
}
