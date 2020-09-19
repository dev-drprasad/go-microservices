package grpc

import (
	"context"
	"gomicroservices/internal/common/genproto/organization"
	orgserv "gomicroservices/internal/organization/service"
	"gomicroservices/internal/util"
)

var service orgserv.Service

func init() {
	service = orgserv.New()
}

type Organization struct {
}

func NewOrganization() *Organization {
	return &Organization{}
}

func (o *Organization) GetBranch(ctx context.Context, in *organization.ByIDRequest) (*organization.Branch, error) {
	b := *service.GetBranch(util.NewContextFromMetadata(ctx), in.ID)
	return &organization.Branch{Name: b.Name, PhoneNumber: b.PhoneNumber}, nil
}
