package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"rpg-backend/database"
	"rpg-backend/handlers"
)

func main() {
	// Inicializar conexão com o banco de dados
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializar handlers
	h := handlers.NewHandlers(db)

	// Configurar roteador
	r := mux.NewRouter()

	// Rotas de autenticação
	r.HandleFunc("/auth/google", h.HandleGoogleAuth).Methods("POST")

	// Rotas de usuário (protegidas)
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(h.AuthMiddleware)
	userRouter.HandleFunc("/me", h.GetCurrentUser).Methods("GET")

	// Rotas de personagens (protegidas)
	charRouter := r.PathPrefix("/characters").Subrouter()
	charRouter.Use(h.AuthMiddleware)
	charRouter.HandleFunc("", h.CreateCharacter).Methods("POST")
	charRouter.HandleFunc("", h.GetUserCharacters).Methods("GET")
	charRouter.HandleFunc("/{id}", h.GetCharacter).Methods("GET")
	charRouter.HandleFunc("/{id}", h.UpdateCharacter).Methods("PUT")
	charRouter.HandleFunc("/{id}", h.DeleteCharacter).Methods("DELETE")

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}