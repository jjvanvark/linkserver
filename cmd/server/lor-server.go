package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jjvanvark/lor-deluxe/internal/graphql"
	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jjvanvark/lor-deluxe/internal/routes"
	"github.com/jjvanvark/lor-deluxe/internal/security"
	"github.com/jjvanvark/lor-deluxe/internal/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	var stop chan os.Signal
	var router chi.Router
	var rts *routes.Routes
	var auth *jwtauth.JWTAuth
	var server *http.Server
	var db models.DataInterface
	var graphqlHandler http.Handler

	// env
	var db_name string = ":memory:"
	var server_addr string = ":9090"

	if err = godotenv.Load(); err == nil {
		db_name = os.Getenv("DB_NAME")
		server_addr = os.Getenv("SERVER_ADDRESS")
	}

	defer shutdown(server, db)

	// database
	if db, err = sqlite.Init(db_name); err != nil {
		fmt.Printf("Database init error :: %v\n", err)
		return
	}

	optionallyAddUser(db)

	// graphql
	if graphqlHandler, err = graphql.InitGraphql(db); err != nil {
		fmt.Printf("Graphql init error :: %v\n", err)
		return
	}

	// router
	auth = jwtauth.New("HS256", []byte(os.Getenv("AUTH_KEY")), nil)
	rts = routes.InitRoutes(db, auth)
	router = chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth))
		r.Use(jwtauth.Authenticator)

		r.Handle("/query", graphqlHandler)
		r.Get("/asset/{cursor}", rts.HandleAsset)
	})

	router.Group(func(r chi.Router) {
		r.Post("/login", rts.HandleLogin)
		r.Get("/login/{email}/{password}", rts.HandleHttpLogin)
		r.Get("/addurl", rts.HandleAddUrl)
	})

	// server
	server = &http.Server{
		Addr:    server_addr,
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			fmt.Printf("Server error :: %v\n", err)
		}
	}()

	// shuting down
	fmt.Printf("Listening on 9090 \n")
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func shutdown(server *http.Server, db models.DataInterface) {

	fmt.Printf("\nShutting down gracefully\n")

	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	if server != nil {
		if err = server.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down server :: %v\n", err)
		}
	}

	if db != nil {
		if err = db.Close(); err != nil {
			fmt.Printf("Error closing db :: %v\n", err)
		}
	}

	fmt.Printf("Shutting down done\n")

}

func optionallyAddUser(d models.DataInterface) error {

	var err error
	var pw []byte
	var key string

	if _, err = d.GetUserByEmail("mail@jvanvark.nl"); err != nil {
		if err == sql.ErrNoRows {
			if pw, err = security.HashPassword("joost"); err != nil {
				fmt.Printf("Password hash error :: %v\n", err)
				return err
			}
			if _, err = d.AddUser("mail@jvanvark.nl", "Joost", pw, func(id int64) (string, error) {
				if key, err = security.GetUserToken(id); err != nil {
					return "", err
				}
				return key, nil
			}); err != nil {
				fmt.Printf("Add user error :: %v\n", err)
				return err
			}
		} else {
			fmt.Printf("Get user by email error :: %v\n", err)
			return err
		}
	}

	return nil

}
