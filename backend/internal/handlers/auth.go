package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/rishavtarway/taskflow/internal/config"
	"github.com/rishavtarway/taskflow/internal/errors"
	"github.com/rishavtarway/taskflow/internal/models"
)

type AuthHandler struct {
	userModel *models.UserModel
	cfg       *config.Config
}

func NewAuthHandler(userModel *models.UserModel, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userModel: userModel,
		cfg:       cfg,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string               `json:"token"`
	User  models.UserWithToken `json:"user"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	fields := make(map[string]string)
	if req.Name == "" {
		fields["name"] = "is required"
	}
	if req.Email == "" {
		fields["email"] = "is required"
	} else if !emailRegex.MatchString(req.Email) {
		fields["email"] = "invalid format"
	}
	if req.Password == "" {
		fields["password"] = "is required"
	} else if len(req.Password) < 8 {
		fields["password"] = "must be at least 8 characters"
	}

	if len(fields) > 0 {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(fields))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	user, err := h.userModel.Create(r.Context(), req.Name, req.Email, string(hashedPassword))
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` ||
			err.Error() == `pq: duplicate key value violates unique constraint "users_email_unique"` {
			errors.WriteError(w, errors.Conflict("email already exists"))
			return
		}
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user.ToUserWithToken()})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	fields := make(map[string]string)
	if req.Email == "" {
		fields["email"] = "is required"
	}
	if req.Password == "" {
		fields["password"] = "is required"
	}

	if len(fields) > 0 {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(fields))
		return
	}

	user, err := h.userModel.GetByEmail(r.Context(), req.Email)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	if user == nil {
		bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		errors.WriteError(w, errors.Unauthorized())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		errors.WriteError(w, errors.Unauthorized())
		return
	}

	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user.ToUserWithToken()})
}

func (h *AuthHandler) generateToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"iat":   now.Unix(),
		"exp":   now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
