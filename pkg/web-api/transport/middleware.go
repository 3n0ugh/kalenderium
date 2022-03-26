package transport

import (
	"context"
	"fmt"
	contx "github.com/3n0ugh/kalenderium/internal/context"
	errs "github.com/3n0ugh/kalenderium/internal/err"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	repo "github.com/3n0ugh/kalenderium/pkg/account/repository"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the clients' IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {

			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract the client's IP address from the request.
		ip := realip.FromRequest(r)

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Limit(2),
					4,
				),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			errs.RateLimitExceededResponse(w)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = contx.SetUser(r, repo.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.Write([]byte(fmt.Sprintf("%s\n", headerParts)))
			errs.InvalidAuthenticationTokenResponse(w)
			return
		}

		tkn := headerParts[1]

		v := validator.New()

		if token.ValidateTokenPlaintext(v, tkn); !v.Valid() {
			errs.InvalidAuthenticationTokenResponse(w)
			return
		}

		var grpcAddr = net.JoinHostPort("localhost", "8083") // account service
		conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewAccountClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		pToken := &pb.Token{
			PlaintText: tkn,
			Hash:       nil,
			UserId:     0,
			Expiry:     nil,
			Scope:      token.ScopeAuthentication,
		}

		resp, err := client.IsAuth(ctx, &pb.IsAuthRequest{Token: pToken})
		if err != nil {
			errs.InvalidAuthenticationTokenResponse(w)
			return
		}

		usr := &repo.User{
			UserID:       resp.Token.UserId,
			Email:        "",
			Password:     "",
			PasswordHash: nil,
		}
		r = contx.SetUser(r, usr)

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func requireAuthenticatedUser(next *httpTransport.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := contx.GetUser(r)

		if user.IsAnonymous() {
			errs.AuthenticationRequiredResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := w.Header().Get("Origin")

		if origin == "localhost:8080" {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if r.Method == http.MethodOptions &&
				w.Header().Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods",
					"OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers",
					"Authorization, Content-Type")

				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				errs.ServerErrorResponse(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
