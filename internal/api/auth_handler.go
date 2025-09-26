package api

import (
	"encoding/json"
	"net/http"

	"github.com/deepraj02/go-postgres-starter/internal/middleware"
	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/auth"
	utils "github.com/deepraj02/go-postgres-starter/internal/utils/json"
	"github.com/deepraj02/go-postgres-starter/internal/utils/logger"
)

type AuthHandler struct {
	authStore store.AuthStore
	logger    *logger.Logger
}

func NewAuthHandler(authStore store.AuthStore, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authStore: authStore,
		logger:    logger,
	}
}

func (app *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req store.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Username, email, and password are required",
		})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		app.logger.Error("Failed to hash password", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to process password",
		})
		return
	}

	user := &store.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	if req.Bio != "" {
		user.Bio = &req.Bio
	}

	if err := app.authStore.CreateUser(user); err != nil {
		app.logger.Error("Failed to create user", err)
		utils.WriteJson(w, http.StatusConflict, utils.Envelope{
			"error": "Username or email already exists",
		})
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		app.logger.Error("Failed to generate access token", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate authentication token",
		})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		app.logger.Error("Failed to generate refresh token", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate authentication token",
		})
		return
	}

	response := store.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
		User:         *user,
	}

	app.logger.Info("User registered successfully: %s", user.Username)
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"data": response})
}

func (app *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req store.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Username and password are required",
		})
		return
	}

	user, err := app.authStore.GetUserByUsername(req.Username)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
			"error": "Invalid credentials",
		})
		return
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
			"error": "Invalid credentials",
		})
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		app.logger.Error("Failed to generate access token", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate authentication token",
		})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		app.logger.Error("Failed to generate refresh token", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate authentication token",
		})
		return
	}

	response := store.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
		User:         *user,
	}

	app.logger.Info("User logged in successfully: %s", user.Username)
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"data": response})
}

func (app *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
			"error": "Unauthorized",
		})
		return
	}

	user, err := app.authStore.GetUserByID(claims.UserID)
	if err != nil {
		app.logger.Error("Failed to get user profile", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to get user profile",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"data": user})
}
