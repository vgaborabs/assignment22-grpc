package user

import (
	"context"
	"github.com/vgaborabs/assignment22-grpc/internal/db"
	pb "github.com/vgaborabs/assignment22-grpc/proto"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	repo db.UserRepository
}

func NewUserService(repo db.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUser(ctx context.Context, request *pb.IdRequest) (*pb.UserModel, error) {
	u, err := s.repo.GetUserById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return convertUser(u), nil
}

func (s *Service) GetUsers(ctx context.Context, request *pb.MultipleIdRequest) (*pb.UserModels, error) {
	users, err := s.repo.GetUsersByIds(ctx, request.Ids)
	if err != nil {
		return nil, err
	}
	return convertUsers(users), nil
}

func (s *Service) SearchUsers(ctx context.Context, request *pb.SearchCriteria) (*pb.UserModels, error) {
	users, err := s.repo.SearchUsers(ctx, convertSearchCriteria(request))
	if err != nil {
		return nil, err
	}
	return convertUsers(users), nil
}

func convertSearchCriteria(r *pb.SearchCriteria) db.SearchCriteria {
	sc := db.SearchCriteria{
		Field: r.Field,
		Value: r.Value,
	}
	if r.MatchMode != nil {
		dbmm := db.MatchMode(*r.MatchMode)
		sc.MatchMode = &dbmm
	}
	return sc
}

func convertUsers(u []db.User) *pb.UserModels {
	pbu := &pb.UserModels{
		Users: make([]*pb.UserModel, len(u)),
	}
	for i, user := range u {
		pbu.Users[i] = convertUser(user)
	}
	return pbu
}

func convertUser(u db.User) *pb.UserModel {
	return &pb.UserModel{
		Id:      u.Id,
		Fname:   u.FirstName,
		City:    u.City,
		Phone:   u.PhoneNumber,
		Height:  u.Height,
		Married: u.Married,
	}
}
