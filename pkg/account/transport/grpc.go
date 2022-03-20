package transport

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account/endpoints"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	isAuth grpcTransport.Handler
	signUp grpcTransport.Handler
	login  grpcTransport.Handler
	logout grpcTransport.Handler
}

func NewGRPCServer(ep endpoints.Set) pb.AccountServer {
	return &gRPCServer{
		isAuth: grpcTransport.NewServer(
			ep.IsAuthEndpoint,
			decodeIsAuthRequest,
			encodeIsAuthResponse),
		signUp: grpcTransport.NewServer(
			ep.SignUpEndpoint,
			decodeSignUpRequest,
			encodeSignUpResponse),
		login: grpcTransport.NewServer(
			ep.LoginEndpoint,
			decodeLoginRequest,
			encodeLoginResponse),
		logout: grpcTransport.NewServer(
			ep.LogoutEndpoint,
			decodeLogoutRequest,
			encodeLogoutResponse),
	}
}

func (g *gRPCServer) IsAuth(ctx context.Context, r *pb.IsAuthRequest) (*pb.IsAuthReply, error) {
	_, resp, err := g.isAuth.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.IsAuthReply), nil
}

func (g *gRPCServer) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.SignUpReply, error) {
	_, resp, err := g.signUp.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SignUpReply), nil
}

func (g *gRPCServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginReply), nil
}

func (g *gRPCServer) Logout(ctx context.Context, r *pb.LogoutRequest) (*pb.LogoutReply, error) {
	_, resp, err := g.logout.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LogoutReply), nil
}

// decodeIsAuthRequest extracts a user-domain request object from a gRPC request
func decodeIsAuthRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.IsAuthRequest)
	return endpoints.IsAuthRequest{Token: request.Token}, nil
}

// encodeIsAuthResponse encodes the passed response object to the gRPC response message.
func encodeIsAuthResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(*pb.IsAuthReply)
	return endpoints.IsAuthResponse{Err: reply.Err}, nil
}

// decodeSignUpRequest extracts a user-domain request object from a gRPC request
func decodeSignUpRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.SignUpRequest)
	user := repository.User{
		Email:    request.User.Email,
		Password: request.User.Password,
	}
	return endpoints.SignUpRequest{User: user}, nil
}

// encodeSignUpResponse encodes the passed response object to the gRPC response message.
func encodeSignUpResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(*pb.SignUpReply)
	return endpoints.SignUpResponse{Token: reply.Token, Err: reply.Err}, nil
}

// decodeLoginRequest extracts a user-domain request object from a gRPC request
func decodeLoginRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.LoginRequest)
	user := repository.User{
		Email:    request.User.Email,
		Password: request.User.Password,
	}
	return endpoints.LoginRequest{User: user}, nil
}

// encodeLoginResponse encodes the passed response object to the gRPC response message.
func encodeLoginResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(*pb.LoginReply)
	return endpoints.LoginResponse{Token: reply.Token, Err: reply.Err}, nil
}

// decodeLogoutRequest extracts a user-domain request object from a gRPC request
func decodeLogoutRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.LogoutRequest)
	return endpoints.LogoutRequest{Token: request.Token}, nil
}

// encodeLogoutResponse encodes the passed response object to the gRPC response message.
func encodeLogoutResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(*pb.LogoutReply)
	return endpoints.LogoutResponse{Err: reply.Err}, nil
}