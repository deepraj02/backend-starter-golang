package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/deepraj02/go-postgres-starter/internal/middleware"
	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/auth"
	utils "github.com/deepraj02/go-postgres-starter/internal/utils/json"
	"github.com/deepraj02/go-postgres-starter/internal/utils/logger"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandler struct {
	authStore    store.AuthStore
	logger       *logger.Logger
	cacheStore   store.CacheStore
	emailService store.EmailService
	sessionStore *sessions.CookieStore
}

type contextKey string

const providerKey = contextKey("provider")

func NewAuthHandler(authStore store.AuthStore, logger *logger.Logger, cacheStore store.CacheStore, emailService store.EmailService) *AuthHandler {
	
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, 
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
	}

	
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
		),
	)

	gothic.Store = sessionStore

	return &AuthHandler{
		authStore:    authStore,
		logger:       logger,
		cacheStore:   cacheStore,
		emailService: emailService,
		sessionStore: sessionStore,
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

	
	existingUser, err := app.authStore.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		utils.WriteJson(w, http.StatusConflict, utils.Envelope{
			"error": fmt.Sprintf("An account with this email already exists using %s authentication", existingUser.Provider),
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
		Provider:     "email", 
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
		ExpiresIn:    24 * 60 * 60, 
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

	
	if user.Provider != "email" {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
			"error": fmt.Sprintf("This account was created using %s  Authentication. Please use that method to login.", user.Provider),
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


func (app *AuthHandler) BeginOAuth(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		provider = "google" 
	}
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))
	
	session, _ := app.sessionStore.Get(r, "goth-session")
	session.Values["provider"] = provider
	session.Save(r, w)

	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))
	gothic.BeginAuthHandler(w, r)
}

func (app *AuthHandler) CompleteOAuth(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.logger.Error("OAuth authentication failed", err)
		http.Redirect(w, r, "/login.html?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	
	existingUser, err := app.authStore.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil && existingUser.Provider != strings.ToLower(user.Provider) {
		app.logger.Error("User exists with different provider", fmt.Errorf("email %s already registered with %s", user.Email, existingUser.Provider))
		http.Redirect(w, r, fmt.Sprintf("/login.html?error=email_exists&provider=%s", existingUser.Provider), http.StatusTemporaryRedirect)
		return
	}

	var dbUser *store.User

	if existingUser != nil {
		
		dbUser = existingUser
	} else {
		
		newUser := &store.User{
			Username: user.Email, 
			Email:    user.Email,
			Provider: strings.ToLower(user.Provider),
		}

		if user.Name != "" {
			newUser.Bio = &user.Name
		}

		if err := app.authStore.CreateOAuthUser(newUser); err != nil {
			app.logger.Error("Failed to create OAuth user", err)
			http.Redirect(w, r, "/login.html?error=registration_failed", http.StatusTemporaryRedirect)
			return
		}
		dbUser = newUser
	}

	
	accessToken, err := auth.GenerateJWT(dbUser.ID, dbUser.Email)
	if err != nil {
		app.logger.Error("Failed to generate access token for OAuth user", err)
		http.Redirect(w, r, "/login.html?error=token_generation_failed", http.StatusTemporaryRedirect)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(dbUser.ID, dbUser.Email)
	if err != nil {
		app.logger.Error("Failed to generate refresh token for OAuth user", err)
		http.Redirect(w, r, "/login.html?error=token_generation_failed", http.StatusTemporaryRedirect)
		return
	}

	redirectURL := fmt.Sprintf("/dashboard.html?access_token=%s&refresh_token=%s&user=%s",
		accessToken, refreshToken,
		url.QueryEscape(fmt.Sprintf(`{"id":%d,"username":"%s","email":"%s","provider":"%s"}`,
			dbUser.ID, dbUser.Username, dbUser.Email, dbUser.Provider)))

	app.logger.Info("OAuth user logged in successfully: %s", dbUser.Email)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
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

	
	user, err := h.authStore.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{
			"error": "Email not found",
		})
		return
	}

	if user.Provider != "email" {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{
			"error": fmt.Sprintf("This account uses %s authentication. Password reset is not available.", user.Provider),
		})
		return
	}

	
	code, err := store.GenerateResetCode()
	if err != nil {
		h.logger.Error("Failed to generate reset code", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to generate reset code",
		})
		return
	}

	
	if err := h.cacheStore.SetResetCode(req.Email, code, 15*time.Minute); err != nil {
		h.logger.Error("Failed to store reset code", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to process request",
		})
		return
	}

	email := os.Getenv("EMAIL_FROM")
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

	
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to hash new password", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to process new password",
		})
		return
	}

	if err := h.authStore.UpdatePassword(req.Email, hashedPassword); err != nil {
		h.logger.Error("Failed to update password in database", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to update password",
		})
		return
	}

	
	if err := h.cacheStore.DeleteResetCode(req.Email); err != nil {
		h.logger.Error("Failed to delete reset code", err)
	}

	h.logger.Info("Password reset successfully for email: %s", req.Email)
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "Password reset successfully",
	})
}
