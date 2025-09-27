package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/deepraj02/go-postgres-starter/internal/middleware"
	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/auth"
	utils "github.com/deepraj02/go-postgres-starter/internal/utils/json"
	"github.com/deepraj02/go-postgres-starter/internal/utils/logger"
)

type AuthHandler struct {
	authStore    store.AuthStore
	logger       *logger.Logger
	cacheStore   store.CacheStore
	emailService store.EmailService
}

func NewAuthHandler(authStore store.AuthStore, logger *logger.Logger, cacheStore store.CacheStore, emailService store.EmailService) *AuthHandler {
	return &AuthHandler{
		authStore:    authStore,
		logger:       logger,
		cacheStore:   cacheStore,
		emailService: emailService,
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
		ExpiresIn:    24 * 60 * 60,
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

// ----------------- PASSWORD ------------

// ...existing code...
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req store.ForogtPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}

	if req.Email == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Email is required",
		})
		return
	}

	// Check if user exists
	_, err := h.authStore.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{
			"error": "Email not found",
		})
		return
	}

	// Generate 6-digit code
	code, err := store.GenerateResetCode()
	if err != nil {
		h.logger.Error("Failed to generate reset code", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate reset code",
		})
		return
	}

	// Store code in Redis with 15-minute expiration
	if err := h.cacheStore.SetResetCode(req.Email, code, 15*time.Minute); err != nil {
		h.logger.Error("Failed to store reset code", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to process request",
		})
		return
	}

	email  := os.Getenv("EMAIL_FROM")
	emailBody := store.ResendEmailBody{
		From:    email, 
		To:      []string{req.Email},
		Subject: "Password Reset Code",
		HTML:    h.emailService.(*store.ResendEmailService).GenerateResetEmailHTML(code),
	}

	if err := h.emailService.SendResetCode(emailBody); err != nil {
		h.logger.Error("Failed to send reset email", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to send reset email",
		})
		return
	}

	h.logger.Info("Reset code sent to email: %s", req.Email)
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "Reset code sent to your email",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req store.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}

	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Email, code, and new password are required",
		})
		return
	}

	// Verify code from Redis
	storedCode, err := h.cacheStore.GetResetCode(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid or expired reset code",
		})
		return
	}

	if storedCode != req.Code {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid reset code",
		})
		return
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to hash new password", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to process new password",
		})
		return
	}

	// Update password in database - NOW ACTUALLY IMPLEMENTED!
	if err := h.authStore.UpdatePassword(req.Email, hashedPassword); err != nil {
		h.logger.Error("Failed to update password in database", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to update password",
		})
		return
	}

	// Delete the used reset code
	if err := h.cacheStore.DeleteResetCode(req.Email); err != nil {
		h.logger.Error("Failed to delete reset code", err)
	}

	h.logger.Info("Password reset successfully for email: %s", req.Email)
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "Password reset successfully",
	})
}
