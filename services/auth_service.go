package services

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/config"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/models"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/repositories"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}

func (s *AuthService) VerifyFirebaseToken(firebaseToken string) (string, *models.User, error){
	token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), firebaseToken)
	if err != nil {
		return "", nil, errors.New("firebase token tidak valid atau kadaluarsa")
	}

	emailVerified, _ := token.Claims["email_verified"].(bool)
	if !emailVerified {
		return "", nil, errors.New("EMAIL_NOT_VERIFIED")
	}


	uid := token.UID
	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	user, err := s.userRepo.FindByFirebaseUID(uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		now := time.Now().Unix()
		user = &models.User{
			FirebaseUID:   uid,
			Email:         email,
			Name:          name,
			Role: 		   "user",
			EmailVerified: true,
			LastLoginAt:   &now,
		}
		if err := s.userRepo.Create(user); err != nil {
			return "", nil, errors.New("gagal bikin user baru")
		}
	} else if err != nil {
		return "", nil, errors.New("gagal mengambil data user")
	} else {
		now := time.Now().Unix()
		user.LastLoginAt = &now
		user.EmailVerified = true
		s.userRepo.Update(user)
	}

	jwtToken, err := s.generateJWT(user)
	if err != nil {
		return "", nil, errors.New("gagal membuat token")
	}

	return jwtToken, user, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if expireHours == 0 {
		expireHours = 24
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"firebase_uid": user.FirebaseUID,
		"email": user.Email,
		"name": user.Name,
		"role": user.Role,
		"email_verified": user.EmailVerified,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour *time.Duration(expireHours)).Unix(),
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}