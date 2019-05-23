package main

import (
	"context"
	pb "github.com/bukhavtsov/gateway-jwt-auth/pkg/proto"
	auth "github.com/bukhavtsov/restful-app/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	client auth.AuthClient
}

func (s *server) ReadAllDevelopers(ctx context.Context, req *pb.ReadAllDevelopersRequest) (resp *pb.ReadAllDevelopersResponse, err error) {
	authReq := auth.ReadAllDevelopersRequest{}
	authResp, err := s.client.ReadAllDevelopers(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	var respDevelopers []*pb.Developer
	for _, currentDev := range authResp.Developers {
		developer := &pb.Developer{
			Id:           currentDev.Id,
			Name:         currentDev.Name,
			Age:          currentDev.Age,
			PrimarySkill: currentDev.PrimarySkill,
		}
		respDevelopers = append(respDevelopers, developer)
	}
	return &pb.ReadAllDevelopersResponse{Developers: respDevelopers}, nil
}

func (s *server) CreateDeveloper(ctx context.Context, req *pb.CreateDeveloperRequest) (*pb.CreateDeveloperResponse, error) {
	developer := auth.Developer{
		Id:           req.Developer.Id,
		Name:         req.Developer.Name,
		Age:          req.Developer.Age,
		PrimarySkill: req.Developer.PrimarySkill,
	}
	authReq := auth.CreateDeveloperRequest{
		Developer: &developer,
	}
	authResp, err := s.client.CreateDeveloper(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	return &pb.CreateDeveloperResponse{Id: authResp.Id}, nil
}

func (s *server) ReadDeveloper(ctx context.Context, req *pb.ReadDeveloperRequest) (*pb.ReadDeveloperResponse, error) {
	authReq := auth.ReadDeveloperRequest{Id: req.Id}
	authResp, err := s.client.ReadDeveloper(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	respDeveloper := &pb.Developer{
		Id:           authResp.Developer.Id,
		Name:         authResp.Developer.Name,
		Age:          authResp.Developer.Age,
		PrimarySkill: authResp.Developer.PrimarySkill,
	}
	return &pb.ReadDeveloperResponse{Developer: respDeveloper}, nil
}

func (s *server) UpdateDeveloper(ctx context.Context, req *pb.UpdateDeveloperRequest) (*pb.UpdateDeveloperResponse, error) {
	developer := auth.Developer{
		Id:           req.Developer.Id,
		Name:         req.Developer.Name,
		Age:          req.Developer.Age,
		PrimarySkill: req.Developer.PrimarySkill,
	}
	authReq := auth.UpdateDeveloperRequest{
		Developer: &developer,
		Id:        req.Id,
	}
	authResp, err := s.client.UpdateDeveloper(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	respDeveloper := &pb.Developer{
		Id:           authResp.Developer.Id,
		Name:         authResp.Developer.Name,
		Age:          authResp.Developer.Age,
		PrimarySkill: authResp.Developer.PrimarySkill,
	}
	return &pb.UpdateDeveloperResponse{Developer: respDeveloper}, nil
}

func (s *server) DeleteDeveloper(ctx context.Context, req *pb.DeleteDeveloperRequest) (*pb.DeleteDeveloperResponse, error) {
	authReq := auth.DeleteDeveloperRequest{
		Id: req.Id,
	}
	_, err := s.client.DeleteDeveloper(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	return &pb.DeleteDeveloperResponse{}, nil
}

func (s *server) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	authReq := auth.SignInRequest{
		Login:    req.Login,
		Password: req.Password,
	}
	authResp, err := s.client.SignIn(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	return &pb.SignInResponse{RefreshToken: authResp.RefreshToken, AccessToken: authResp.AccessToken}, nil
}

func (s *server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	authReq := auth.SignUpRequest{
		Login:    req.Login,
		Password: req.Password,
	}
	authResp, err := s.client.SignUp(ctx, &authReq)
	if err != nil {
		log.Println(err)
	}
	return &pb.SignUpResponse{RefreshToken: authResp.RefreshToken, AccessToken: authResp.AccessToken}, nil
}

func main() {
	listener, err := net.Listen("tcp", "2020")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	cc, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := auth.NewAuthClient(cc)
	pb.RegisterRestAppServer(srv, &server{client})
	reflection.Register(srv)
	if e := srv.Serve(listener); e != nil {
		log.Fatal(err)
	}
}
