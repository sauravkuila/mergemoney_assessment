package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	v1 "github.com/sauravkuila/mergemoney_assessment/pkg/service/v1"
)

type ServiceItf interface {
	Health(c *gin.Context)
	GetV1Object() v1.V1Group
}

type serviceSt struct {
	v1 v1.V1Group
}

func (obj *serviceSt) GetV1Object() v1.V1Group {
	return obj.v1
}

func GetServiceItf(dao dao.RepositoryItf) ServiceItf {
	return &serviceSt{
		v1: v1.GetV1Group(dao),
	}
}
