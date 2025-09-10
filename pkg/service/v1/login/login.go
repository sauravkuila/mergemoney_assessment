package login

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
)

type LoginItf interface {
}

type loginSt struct {
	DB dao.RepositoryItf
	// utils beutils.UtilsItf
}

func GetLoginItf(dao dao.RepositoryItf) LoginItf {
	// cfg := beutils.UtilsConfig{
	// 	RestSSlEnabled: false,
	// }
	return &loginSt{
		DB: dao,
		// utils: beutils.GetUtilsObj(cfg),
	}
}
