package auth

import (
	"fmt"
	"innotech/pkg/keycloak"
	"log/slog"
	"strings"
)

//TODO: пересмотрель связку с user из pg таблиц

// Service provides authentication and user management services
type Service struct {
	keycloakClient *keycloak.Client
	logger         *slog.Logger
}

// NewService creates a new auth service
func NewService(keycloakClient *keycloak.Client, logger *slog.Logger) *Service {
	return &Service{
		keycloakClient: keycloakClient,
		logger:         logger,
	}
}

// CreateUser creates a new user in Keycloak (admin only)
func (s *Service) CreateUser(req CreateUserRequest) (*CreateUserResponse, error) {
	s.logger.Info("creating new user",
		"email", req.Email,
	)

	existingUser, err := s.keycloakClient.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	username := strings.Split(req.Email, "@")[0]

	createReq := keycloak.CreateUserRequest{
		Username:      username,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Enabled:       true,
		EmailVerified: false,
		Credentials: []keycloak.CredentialRepresentation{
			{
				Type:      "password",
				Value:     req.Password,
				Temporary: true,
			},
		},
	}

	userID, err := s.keycloakClient.CreateUser(createReq)
	if err != nil {
		s.logger.Error("failed to create user in Keycloak",
			"email", req.Email,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("user created successfully",
		"userId", userID,
		"email", req.Email,
	)

	return &CreateUserResponse{
		UserID:   userID,
		Email:    req.Email,
		Username: username,
		Message:  "User created successfully. Password is temporary and must be changed on first login.",
	}, nil
}

func (s *Service) ChangePassword(userID string, req ChangePasswordRequest) error {
	s.logger.Info("changing password",
		"userId", userID,
	)

	err := s.keycloakClient.UpdatePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		s.logger.Error("failed to change password",
			"userId", userID,
			"error", err,
		)
		return fmt.Errorf("failed to change password: %w", err)
	}

	s.logger.Info("password changed successfully",
		"userId", userID,
	)

	return nil
}

// GetUserInfo retrieves user information
func (s *Service) GetUserInfo(userID string) (*UserInfoResponse, error) {
	user, err := s.keycloakClient.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &UserInfoResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		EmailVerified: user.EmailVerified,
		Enabled:       user.Enabled,
	}, nil
}
