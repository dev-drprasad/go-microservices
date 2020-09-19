package organization

import (
	"context"
	"gomicroservices/internal/common/genproto/organization"
	"gomicroservices/internal/organization/model"
	"gomicroservices/internal/util"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GRPCRepo struct{}

func New() GRPCRepo {
	return GRPCRepo{}
}

func (repo GRPCRepo) GetBranch(ctx context.Context, id uint64) *model.Branch {
	log := util.GetLoggerFromContext(ctx)
	serverAddress := "localhost:8082"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "rID", log.Prefix())
	defer cancel()
	log.Infof("Dialing to organization grpc service. address=%s", serverAddress)
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Errorf("Dail failed to organization service. %s", err.Error())
		return nil
	}

	client := organization.NewOrganizationServiceClient(conn)

	b, err := client.GetBranch(ctx, &organization.ByIDRequest{ID: id})
	if err != nil {
		log.Errorf("Failed to read branch from grpc client: %s id=%d", err, id)
		return nil
	}

	return &model.Branch{
		Name:        b.Name,
		PhoneNumber: b.PhoneNumber,
	}
}
