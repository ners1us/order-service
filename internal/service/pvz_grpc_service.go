package service

import (
	"context"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/ners1us/order-service/pkg/generated/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZGrpcService struct {
	proto.UnimplementedPVZServiceServer
	pvzRepository repository.PVZRepository
}

func NewPVZGrpcService(pvzRepository repository.PVZRepository) *PVZGrpcService {
	return &PVZGrpcService{
		pvzRepository: pvzRepository,
	}
}

func (pgs *PVZGrpcService) GetPVZList(_ context.Context, _ *proto.GetPVZListRequest) (*proto.GetPVZListResponse, error) {
	pvzs, err := pgs.pvzRepository.GetAllPVZs()
	if err != nil {
		return nil, err
	}

	response := &proto.GetPVZListResponse{
		Pvzs: make([]*proto.PVZ, 0, len(pvzs)),
	}

	for _, pvz := range pvzs {
		protoPVZ := &proto.PVZ{
			Id:               pvz.ID,
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             pvz.City,
		}
		response.Pvzs = append(response.Pvzs, protoPVZ)
	}

	return response, nil
}
