package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/pkg/config"
	"techmind/schema/ent"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt" // TODO: argon
)

type authService struct {
	userRepo repo.UserRepository
	config   *config.Config
}

func NewService(userRepo repo.UserRepository, config *config.Config) service.AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (token string, expiresAt time.Time, err error) {
	// Получаем пользователя по email
	fmt.Println(email, password)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("user not found: %w", err)
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("invalid password")
	}

	// Генерируем JWT токен
	expiresAt = time.Now().Add(72 * time.Hour) // TODO: make configurable
	token, err = s.generateToken(user.ID, expiresAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, expiresAt, nil
}

func (s *authService) Register(ctx context.Context, name, email, password string) (token string, user *ent.User, err error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user, err = s.userRepo.Create(ctx, name, email, string(hashedPassword))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем JWT токен для автоматического входа
	expiresAt := time.Now().Add(72 * time.Hour)
	token, err = s.generateToken(user.ID, expiresAt)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (userID uuid.UUID, err error) {
	// Парсим и валидируем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return uuid.Nil, errors.New("token expired")
			}
		}

		// Извлекаем user_id
		if userIDStr, ok := claims["user_id"].(string); ok {
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return uuid.Nil, fmt.Errorf("invalid user_id in token: %w", err)
			}
			return userID, nil
		}

		return uuid.Nil, errors.New("user_id not found in token")
	}

	return uuid.Nil, errors.New("invalid token")
}

func (s *authService) generateToken(userID uuid.UUID, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
