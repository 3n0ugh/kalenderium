package account

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCServer struct {
	isAuth        grpcTransport.Handler
	signUp        grpcTransport.Handler
	login         grpcTransport.Handler
	logout        grpcTransport.Handler
	serviceStatus grpcTransport.Handler
}

func NewGRPCServer(ep Set) pb.AccountServer {
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
		serviceStatus: grpcTransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeServiceStatusRequest,
			encodeServiceStatusResponse),
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

func (g *gRPCServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ServiceStatusReply), nil
}

// decodeIsAuthRequest extracts a user-domain request object from a gRPC request
func decodeIsAuthRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.IsAuthRequest)

	sessionToken := token.Token{
		PlainText: request.Token.PlaintText,
		Hash:      request.Token.Hash,
		UserID:    request.Token.UserId,
		Expiry:    request.Token.Expiry.AsTime(),
		Scope:     request.Token.Scope,
	}
	return IsAuthRequest{Token: sessionToken}, nil
}

// encodeIsAuthResponse encodes the passed response object to the gRPC response message.
func encodeIsAuthResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(IsAuthResponse)

	tkn := &pb.Token{
		PlaintText: reply.Token.PlainText,
		Hash:       reply.Token.Hash,
		UserId:     reply.Token.UserID,
		Expiry:     timestamppb.New(reply.Token.Expiry),
		Scope:      reply.Token.Scope,
	}
	return &pb.IsAuthReply{Token: tkn, Err: reply.Err}, nil
}

// decodeSignUpRequest extracts a user-domain request object from a gRPC request
func decodeSignUpRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.SignUpRequest)
	user := repository.User{
		Email:    request.User.Email,
		Password: request.User.Password,
	}
	return SignUpRequest{User: user}, nil
}

// encodeSignUpResponse encodes the passed response object to the gRPC response message.
func encodeSignUpResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(SignUpResponse)

	sessionToken := &pb.Token{
		PlaintText: reply.Token.PlainText,
		Hash:       reply.Token.Hash,
		UserId:     reply.Token.UserID,
		Expiry:     timestamppb.New(reply.Token.Expiry),
		Scope:      reply.Token.Scope,
	}
	return &pb.SignUpReply{UserId: reply.UserId, Token: sessionToken, Err: reply.Err}, nil
}

// decodeLoginRequest extracts a user-domain request object from a gRPC request
func decodeLoginRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.LoginRequest)
	user := repository.User{
		Email:    request.User.Email,
		Password: request.User.Password,
	}
	return LoginRequest{User: user}, nil
}

// encodeLoginResponse encodes the passed response object to the gRPC response message.
func encodeLoginResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(LoginResponse)
	sessionToken := &pb.Token{
		PlaintText: reply.Token.PlainText,
		Hash:       reply.Token.Hash,
		UserId:     reply.Token.UserID,
		Expiry:     timestamppb.New(reply.Token.Expiry),
		Scope:      reply.Token.Scope,
	}

	return &pb.LoginReply{UserId: reply.UserId, Token: sessionToken, Err: reply.Err}, nil
}

// decodeLogoutRequest extracts a user-domain request object from a gRPC request
func decodeLogoutRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.LogoutRequest)

	sessionToken := token.Token{
		PlainText: request.Token.PlaintText,
		Hash:      request.Token.Hash,
		UserID:    request.Token.UserId,
		Expiry:    request.Token.Expiry.AsTime(),
		Scope:     request.Token.Scope,
	}

	return LogoutRequest{Token: sessionToken}, nil
}

// encodeLogoutResponse encodes the passed response object to the gRPC response message.
func encodeLogoutResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(LogoutResponse)
	return &pb.LogoutReply{Err: reply.Err}, nil
}

// decodeServiceStatusRequest extracts a user-domain request object from a gRPC request
func decodeServiceStatusRequest(_ context.Context, req interface{}) (interface{}, error) {
	_ = req.(*pb.ServiceStatusRequest)
	return ServiceStatusRequest{}, nil
}

// encodeServiceStatusResponse encodes the passed response object to the gRPC response message.
func encodeServiceStatusResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(ServiceStatusResponse)
	return &pb.ServiceStatusReply{Code: int32(reply.Code), Err: reply.Err}, nil
}
