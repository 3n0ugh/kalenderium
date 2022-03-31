package account

import (
	"context"
	"net"
	"testing"

	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	mockRepo "github.com/3n0ugh/kalenderium/pkg/account/repository/mock"
	mockStore "github.com/3n0ugh/kalenderium/pkg/account/store/mock"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Still getting network behavior, but over an in-memory connection without using OS-level resources
func server(ctx context.Context) (pb.AccountClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	repo := mockRepo.NewAccountRepository()
	redis := mockStore.CustomRedisStore(ctx)
	svc := NewService(repo, redis)
	ep := New(svc)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	pb.RegisterAccountServer(baseServer, NewGRPCServer(ep))
	go func() {
		if err := baseServer.Serve(listener); err != nil {
			logger.Log("err", err)
		}
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	closer := func() {
		listener.Close()
		baseServer.Stop()
	}

	client := pb.NewAccountClient(conn)

	return client, closer
}

// Custom struct comparer for Login handler test
func deepEqualLogin(x, y *pb.LoginReply) bool {
	if x.Err != y.Err {
		return false
	}
	if x.UserId != y.UserId {
		return false
	}
	if x.Token.Expiry.AsTime() != y.Token.Expiry.AsTime() {
		return false
	}
	if string(x.Token.Hash) != string(y.Token.Hash) {
		return false
	}
	if x.Token.PlaintText != y.Token.PlaintText {
		return false
	}
	if x.Token.Scope != y.Token.Scope {
		return false
	}
	return true
}

// Custom struct comparer for Signup handler test
func deepEqualSignup(x, y *pb.SignUpReply) bool {
	if x.Err != y.Err {
		return false
	}
	if x.UserId != y.UserId {
		return false
	}
	if x.Token.Expiry.AsTime() != y.Token.Expiry.AsTime() {
		return false
	}
	if string(x.Token.Hash) != string(y.Token.Hash) {
		return false
	}
	if x.Token.PlaintText != y.Token.PlaintText {
		return false
	}
	if x.Token.Scope != y.Token.Scope {
		return false
	}
	return true
}

// Custom struct comparer for IsAuth handler test
func deepEqualIsAuth(x, y *pb.IsAuthReply) bool {
	if x.Err != y.Err {
		return false
	}
	if x.Token.Expiry.AsTime() != y.Token.Expiry.AsTime() {
		return false
	}
	if string(x.Token.Hash) != string(y.Token.Hash) {
		return false
	}
	if x.Token.PlaintText != y.Token.PlaintText {
		return false
	}
	if x.Token.Scope != y.Token.Scope {
		return false
	}
	return true
}

// Custom struct comparer for IsAuth handler test
func deepEqualLogout(x, y *pb.LogoutReply) bool {
	return x.Err == y.Err
}

func TestAccountService_ServiceStatus(t *testing.T) {
	expected := struct {
		Code int32
		Err  string
	}{
		Code: 200,
		Err:  "",
	}

	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	out, err := client.ServiceStatus(ctx, &pb.ServiceStatusRequest{})

	if err == nil {
		if out.Code != expected.Code {
			t.Errorf("Code -> Want: %d;Got: %d", expected.Code, out.Code)
		}
	}
}

func TestAccountService_SignUp(t *testing.T) {

	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.SignUpReply
		err error
	}

	tests := map[string]struct {
		in       *pb.SignUpRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    "test@test.com",
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					UserId: 23,
					Token: &pb.Token{
						PlaintText: mockStore.Token.PlainText,
						Hash:       mockStore.Token.Hash,
						UserId:     23,
						Expiry:     timestamppb.New(mockStore.Token.Expiry),
						Scope:      mockStore.Token.Scope,
					},
					Err: "",
				},
				err: nil,
			},
		},
		"Empty_Email": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    "",
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					Err: "rpc error: code = Unknown desc = failed user data validation: map[email:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = failed user data validation: map[email:must be provided]"),
			},
		},
		"Empty_Password": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "",
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					Err: "rpc error: code = Unknown desc = failed to create new user",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to create new user"),
			},
		},
		"Invalid_Email_Format": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    "test.com",
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					Err: "rpc error: code = Unknown desc = failed user data validation: map[email:must be valid email address]",
				},
				err: errors.New("rpc error: code = Unknown desc = failed user data validation: map[email:must be valid email address]"),
			},
		},
		"Short_Password": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "1234",
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					Err: "rpc error: code = Unknown desc = failed user data validation: map[password:must be at least 8 bytes long]",
				},
				err: errors.New("rpc error: code = Unknown desc = failed user data validation: map[password:must be at least 8 bytes long]"),
			}},
		"Duplicate_Email": {
			in: &pb.SignUpRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.SignUpReply{
					Err: "rpc error: code = Unknown desc = failed to create new user",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to create new user"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.SignUp(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualSignup(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}

func TestAccountService_Login(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.LoginReply
		err error
	}

	tests := map[string]struct {
		in       *pb.LoginRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					UserId: mockRepo.User.UserID,
					Token: &pb.Token{
						PlaintText: mockStore.Token.PlainText,
						Hash:       mockStore.Token.Hash,
						UserId:     mockStore.Token.UserID,
						Expiry:     timestamppb.New(mockStore.Token.Expiry),
						Scope:      mockStore.Token.Scope,
					},
					Err: "",
				},
				err: nil,
			},
		},
		"Empty_Password": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "",
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					Token: &pb.Token{
						Expiry: &timestamppb.Timestamp{
							Seconds: -62135596800,
						},
					},
					Err: "rpc error: code = Unknown desc = wrong password",
				},
				err: errors.New("rpc error: code = Unknown desc = wrong password"),
			},
		},
		"Short_Password": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "test_1",
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					Err: "rpc error: code = Unknown desc =  failed user data validation",
				},
				err: errors.New("rpc error: code = Unknown desc =  failed user data validation"),
			},
		},
		"Long_Password": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttsettesttesttestt",
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					Err: "rpc error: code = Unknown desc =  failed user data validation",
				},
				err: errors.New("rpc error: code = Unknown desc =  failed user data validation"),
			},
		},
		"Wrong_Password": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    mockRepo.User.Email,
					Password: "test1234",
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					UserId: mockRepo.User.UserID,
					Token: &pb.Token{
						PlaintText: mockStore.Token.PlainText,
						Hash:       mockStore.Token.Hash,
						UserId:     mockStore.Token.UserID,
						Expiry:     timestamppb.New(mockStore.Token.Expiry),
						Scope:      mockStore.Token.Scope,
					},
					Err: "rpc error: code = Unknown desc = wrong password",
				},
				err: errors.New("rpc error: code = Unknown desc = wrong password"),
			},
		},
		"Empty_Email": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    "",
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					UserId: mockRepo.User.UserID,
					Token: &pb.Token{
						PlaintText: mockStore.Token.PlainText,
						Hash:       mockStore.Token.Hash,
						UserId:     mockStore.Token.UserID,
						Expiry:     timestamppb.New(mockStore.Token.Expiry),
						Scope:      mockStore.Token.Scope,
					},
					Err: "rpc error: code = Unknown desc =  failed user data validation",
				},
				err: errors.New("rpc error: code = Unknown desc =  failed user data validation"),
			},
		},
		"Wrong_Email_Format": {
			in: &pb.LoginRequest{
				User: &pb.User{
					Email:    "test@test",
					Password: mockRepo.User.Password,
				},
			},
			expected: expectation{
				out: &pb.LoginReply{
					Err: "rpc error: code = Unknown desc = record not found",
				},
				err: errors.New("rpc error: code = Unknown desc = record not found"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Login(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualLogin(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}

func TestAccountService_IsAuth(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.IsAuthReply
		err error
	}

	tests := map[string]struct {
		in       *pb.IsAuthRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.IsAuthRequest{
				Token: &pb.Token{
					PlaintText: mockStore.Token.PlainText,
					UserId:     mockStore.Token.UserID,
					Expiry:     timestamppb.New(mockStore.Token.Expiry),
				},
			},
			expected: expectation{
				out: &pb.IsAuthReply{
					Token: &pb.Token{
						PlaintText: mockStore.Token.PlainText,
						Hash:       mockStore.Token.Hash,
						UserId:     mockStore.Token.UserID,
						Expiry:     timestamppb.New(mockStore.Token.Expiry),
						Scope:      mockStore.Token.Scope,
					},
					Err: "",
				},
				err: nil,
			},
		},
		"Empty_Token": {
			in: &pb.IsAuthRequest{
				Token: &pb.Token{
					PlaintText: "",
					Hash:       nil,
					UserId:     0,
					Expiry:     nil,
					Scope:      "",
				},
			},
			expected: expectation{
				out: &pb.IsAuthReply{
					Token: nil,
					Err:   "rpc error: code = Unknown desc = failed to validate token: map[token:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to validate token: map[token:must be provided]"),
			},
		},
		"Invalid_Token": {
			in: &pb.IsAuthRequest{
				Token: &pb.Token{
					PlaintText: "123",
					Hash:       nil,
					UserId:     0,
					Expiry:     nil,
					Scope:      "",
				},
			},
			expected: expectation{
				out: &pb.IsAuthReply{
					Token: nil,
					Err:   "rpc error: code = Unknown desc = failed to validate token: map[token:must be at least 26 bytes]",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to validate token: map[token:must be at least 26 bytes]"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.IsAuth(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
					t.Errorf("%v", out.Err)
				}
			} else {
				if !deepEqualIsAuth(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}

func TestAccountService_Logout(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.LogoutReply
		err error
	}

	tests := map[string]struct {
		in       *pb.LogoutRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.LogoutRequest{
				Token: &pb.Token{
					PlaintText: mockStore.Token.PlainText,
					Hash:       mockStore.Token.Hash,
					UserId:     mockStore.Token.UserID,
					Expiry:     timestamppb.New(mockStore.Token.Expiry),
					Scope:      mockStore.Token.Scope,
				},
			},
			expected: expectation{
				out: &pb.LogoutReply{
					Err: "",
				},
				err: nil,
			},
		},
		"Invalid_Token": {
			in: &pb.LogoutRequest{
				Token: &pb.Token{
					PlaintText: "a",
					UserId:     mockStore.Token.UserID,
					Expiry:     timestamppb.New(mockStore.Token.Expiry),
				},
			},
			expected: expectation{
				out: &pb.LogoutReply{
					Err: "rpc error: code = Unknown desc = failed to validate token",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to validate token"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Logout(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualLogout(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}
