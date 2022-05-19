package web_api

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/internal/validator"
	pb2 "github.com/3n0ugh/kalenderium/pkg/account/pb"
	repo "github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/calendar/pb"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type webApiService struct {
	calendarClient pb.CalendarClient
	accountClient  pb2.AccountClient
}

func NewWebApiService(calendarClient pb.CalendarClient, accountClient pb2.AccountClient) Service {
	return &webApiService{
		calendarClient: calendarClient,
		accountClient:  accountClient,
	}
}

func (w *webApiService) AddEvent(ctx context.Context, event repository.Event) (string, error) {
	v := validator.New()
	repository.ValidateEvent(v, event)
	if !v.Valid() {
		return "", errors.New(fmt.Sprintf("failed to validate event: %v", v.Errors))
	}

	pEvent := &pb.Event{
		Id:      event.Id.Hex(),
		UserId:  event.UserId,
		Name:    event.Name,
		Details: event.Details,
		Start:   timestamppb.New(event.Start),
		End:     timestamppb.New(event.End),
		Color:   event.Color,
	}

	resp, err := w.calendarClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: pEvent,
	})

	if err != nil {
		return resp.EventId, errors.New(fmt.Sprintf("failed to validate event: %v", v.Errors))
	}

	return resp.EventId, nil
}

func (w *webApiService) ListEvent(ctx context.Context, userId uint64) ([]repository.Event, error) {
	resp, err := w.calendarClient.ListEvent(ctx, &pb.ListEventRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.New("failed to list event")
	}

	var events []repository.Event
	for _, e := range resp.Events {
		objId, err := primitive.ObjectIDFromHex(e.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert object id")
		}
		event := repository.Event{
			Id:      objId,
			UserId:  e.UserId,
			Name:    e.Name,
			Details: e.Details,
			Start:   e.Start.AsTime(),
			End:     e.End.AsTime(),
			Color:   e.Color,
		}
		events = append(events, event)
	}

	return events, nil
}

func (w *webApiService) DeleteEvent(ctx context.Context, eventId string, userId uint64) error {
	_, err := w.calendarClient.DeleteEvent(ctx, &pb.DeleteEventRequest{
		EventId: eventId,
		UserId:  userId,
	})

	if err != nil {
		return errors.New("failed to delete event")
	}

	return nil
}

func (w *webApiService) SignUp(ctx context.Context, user repo.User) (uint64, token.Token, error) {
	err := user.Set(user.Password)
	if err != nil {
		return 0, token.Token{}, errors.Wrap(err, "failed to hash password")
	}

	v := validator.New()
	repo.ValidateUser(v, &user)
	if !v.Valid() {
		return 0, token.Token{}, errors.New(fmt.Sprintf("failed to user: %v", v.Errors))
	}

	usr := pb2.User{
		Email:    user.Email,
		Password: user.Password,
	}

	resp, err := w.accountClient.SignUp(ctx, &pb2.SignUpRequest{User: &usr})
	if err != nil {
		return 0, token.Token{}, errors.Wrap(err, "failed to signup")
	}

	tkn := token.Token{
		PlainText: resp.Token.PlaintText,
		Hash:      resp.Token.Hash,
		UserID:    resp.Token.UserId,
		Expiry:    resp.Token.Expiry.AsTime(),
		Scope:     resp.Token.Scope,
	}
	return resp.UserId, tkn, nil
}

func (w *webApiService) Login(ctx context.Context, user repo.User) (uint64, token.Token, error) {
	err := user.Set(user.Password)
	if err != nil {
		return 0, token.Token{}, errors.Wrap(err, "failed to hash password")
	}

	v := validator.New()
	repo.ValidateUser(v, &user)
	if !v.Valid() {
		return 0, token.Token{}, errors.New(fmt.Sprintf("failed to user: %v", v.Errors))
	}

	usr := pb2.User{
		Email:    user.Email,
		Password: user.Password,
	}

	resp, err := w.accountClient.Login(ctx, &pb2.LoginRequest{
		User: &usr,
	})
	if err != nil {
		return resp.UserId, token.Token{}, errors.New(resp.Err)
	}

	tkn := token.Token{
		PlainText: resp.Token.PlaintText,
		Hash:      resp.Token.Hash,
		UserID:    resp.Token.UserId,
		Expiry:    resp.Token.Expiry.AsTime(),
		Scope:     resp.Token.Scope,
	}

	return resp.UserId, tkn, nil
}
func (w *webApiService) Logout(ctx context.Context, sToken token.Token) error {
	v := validator.New()
	token.ValidateTokenPlaintext(v, sToken.PlainText)
	if !v.Valid() {
		return errors.New(fmt.Sprintf("failed to validate token: %v", v.Errors))
	}

	tkn := &pb2.Token{
		PlaintText: sToken.PlainText,
		Hash:       sToken.Hash,
		UserId:     sToken.UserID,
		Expiry:     timestamppb.New(sToken.Expiry),
		Scope:      sToken.Scope,
	}

	resp, err := w.accountClient.Logout(ctx, &pb2.LogoutRequest{Token: tkn})
	if err != nil {
		return errors.New(resp.Err)
	}

	return nil
}
