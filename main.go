package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/coreycole/go_htmx/db"
	"github.com/coreycole/go_htmx/handle"
	"github.com/coreycole/go_htmx/lib/sb"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handle.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handle.Make(handle.HandleHomeIndex))
	router.Get("/login", handle.Make(handle.HandleLoginIndex))
	router.Get("/login/provider/google", handle.Make(handle.HandleLoginWithGoogle))
	router.Get("/signup", handle.Make(handle.HandleSignupIndex))
	router.Post("/logout", handle.Make(handle.HandleLogoutCreate))
	router.Post("/login", handle.Make(handle.HandleLoginCreate))
	router.Post("/signup", handle.Make(handle.HandleSignupCreate))
	router.Get("/auth/callback", handle.Make(handle.HandleAuthCallback))
	router.Get("/account/setup", handle.Make(handle.HandleAccountSetupIndex))
	router.Post("/account/setup", handle.Make(handle.HandleAccountSetupCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(handle.WithAccountSetup)
		auth.Get("/settings", handle.Make(handle.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handle.Make(handle.HandleSettingsUsernameUpdate))

		auth.Get("/auth/reset-password", handle.Make(handle.HandleResetPasswordIndex))
		auth.Post("/auth/reset-password", handle.Make(handle.HandleResetPasswordCreate))
		auth.Put("/auth/reset-password", handle.Make(handle.HandleResetPasswordUpdate))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
