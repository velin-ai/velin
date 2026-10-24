package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/velin-ai/velin/internal/admin/config"
	apikey "github.com/velin-ai/velin/internal/api/handlers/api_key"
	"github.com/velin-ai/velin/internal/api/handlers/auth"
	"github.com/velin-ai/velin/internal/api/middlewares"
	"github.com/velin-ai/velin/internal/store"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	apiPort := flag.String(
		"api-port",
		":8080",
		"Port for API service",
	)

	adminPort := flag.String(
		"admin-port",
		":9090",
		"Port for admin service",
	)

	setPassword := flag.String(
		"set",
		"",
		"Set or change admin password",
	)

	setDbURL := flag.String(
		"db",
		"",
		"Set or change database url",
	)

	setJWT_Secret := flag.String(
		"jwt",
		"",
		"Set or change jwt secret key",
	)

	flag.Parse()

	if flagWasProvided("set") {
		if *setPassword == "" {
			fmt.Fprintln(os.Stderr, "Error: Password cannot be empty.")
			os.Exit(1)
		}

		if err := configurePassword(*setPassword); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		return
	}

	if flagWasProvided("db") {
		if *setDbURL == "" {
			fmt.Fprintln(os.Stderr, "Error: database URL cannot be empty.")
			os.Exit(1)
		}

		if err := keyring.Set(config.Service, "database_url", *setDbURL); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving database URL:", err)
			os.Exit(1)
		}

		fmt.Println("Database URL saved successfully.")
		return
	}

	if flagWasProvided("jwt") {
		if *setJWT_Secret == "" {
			fmt.Fprintln(os.Stderr, "Error: jwt secret cannot be empty")
			os.Exit(1)
		}

		if err := keyring.Set(config.Service, config.JWTSecret, *setJWT_Secret); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving jwt secret: ", err)
			os.Exit(1)
		}

		fmt.Println("Jwt Secret saved successfully")

		return
	}

	_, err := keyring.Get(config.Service, config.Account)

	if errors.Is(err, keyring.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "Error: Admin password is not configured.")
		fmt.Fprintln(os.Stderr, "Set it with: go run . -set <password>")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing keyring:", err)
		os.Exit(1)
	}

	dbURL, err := keyring.Get(config.Service, "database_url")
	if errors.Is(err, keyring.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "Error: database URL is not configured.")
		fmt.Fprintln(os.Stderr, `Set it with: velin -db="postgres://user:password@localhost:5432/velin?sslmode=disable"`)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing database URL:", err)
		os.Exit(1)
	}

	_, err = keyring.Get(config.Service, config.JWTSecret)

	if errors.Is(err, keyring.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "Error: jwt secret is not configured.")
		fmt.Fprintln(os.Stderr, `Set it with: velin -jwt="super-secret-key"`)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing jwt secret:", err)
		os.Exit(1)
	}

	// connect database
	err = store.ConnectPostgres(dbURL)

	if err != nil {
		log.Fatal(err)
	}

	// create tables
	err = store.CreateTables()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("API port:", *apiPort)
	fmt.Println("Admin port:", *adminPort)

	// api router
	apiRouter := chi.NewRouter()

	apiRouter.Route("/api/v1", func(r chi.Router) {
		// POST /api/v1/auth/register
		r.Post("/auth/register", auth.RegisterHandler)

		// POST /api/v1/auth/login
		r.Post("/auth/login", auth.LoginHandler)

		// POST /api/v1/auth/refresh
		r.Post("/auth/refresh", auth.RefreshHandler)

		// POST /api/v1/api-key
		r.With(
			middlewares.Auth,
		).Post(
			"/api-key", apikey.CreateNewAPIKeyHandler,
		)

		// GET /api/v1/api-keys
		r.With(
			middlewares.Auth,
		).Get(
			"/api-keys", apikey.ListMyAPIKeys,
		)
	})

	apiServer := &http.Server{
		Addr:    fmt.Sprintf("%s%s", "0.0.0.0", *apiPort),
		Handler: apiRouter,
	}

	// start api server
	go func() {
		log.Printf("API server is listening on %s", *apiPort)
		err := apiServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func flagWasProvided(name string) bool {
	provided := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			provided = true
		}
	})

	return provided
}

func configurePassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	oldHash, err := keyring.Get(config.Service, config.Account)

	if errors.Is(err, keyring.ErrNotFound) {
		return savePassword(newPassword)
	}

	if err != nil {
		return fmt.Errorf("accessing keyring: %w", err)
	}

	fmt.Print("Current admin password: ")

	currentPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		return fmt.Errorf("reading current password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(oldHash),
		currentPassword,
	); err != nil {
		return errors.New("incorrect current password")
	}

	return savePassword(newPassword)
}

func savePassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if err := keyring.Set(
		config.Service,
		config.Account,
		string(hash),
	); err != nil {
		return fmt.Errorf("saving password to keyring: %w", err)
	}

	fmt.Println("Admin password saved successfully.")
	return nil
}
