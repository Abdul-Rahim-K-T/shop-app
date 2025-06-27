package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"shop-backend/config"
	"shop-backend/internal/model"
	"shop-backend/internal/repository"
	"shop-backend/pkg/helper"
	jwtutil "shop-backend/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repository.UserRepository
	Cfg      *config.Config
}

func NewAuthService(UserRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		UserRepo: UserRepo,
		Cfg:      cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) error {
	// Trim whitespace from inputs to prevent validation issues
	user.Email = strings.TrimSpace(user.Email)
	user.Phone = strings.TrimSpace(user.Phone)

	log.Printf("Attempting to register user: %s", user.Email)

	// S 1: Check if email already exists
	existing, err := s.UserRepo.FindByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		// Case: Existing user found
		if existing.IsVerified {
			log.Println("User already exists and verified:", existing.Email)
			return errors.New("user already exists")
		}

		// Case 2: User exists but not verified - update name/phone, keep unverified
		log.Println("User exists but not verified, updating basic info")

		existing.FirstName = user.FirstName
		existing.LastName = user.LastName
		existing.Phone = user.Phone
		existing.UpdatedAt = time.Now()

		return s.UserRepo.Update(ctx, existing)
	}

	// Step 2: Handle unexpected DB errors
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		log.Println("Database error during user lookup:", err)
		return err
	}

	// Step 3: New User Registeration (without OTP or password)
	log.Println("Registering new user")

	user.ID = uuid.New().String()
	user.IsVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save user to DB
	err = s.UserRepo.Create(ctx, user)
	if err != nil {
		log.Println("Failed to create user:", err)
		return err
	}

	log.Println("Successfully registerd basic user", user.Email)
	return nil
}

func (s *AuthService) SendOtp(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)
	log.Printf("Sending OTP to : %s", email)

	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Println("User not found:", email)
			return errors.New("user not found")
		}
		log.Println("Error retrieving user:", err)
		return err
	}

	if user.IsVerified {
		log.Println("User already verified:", email)
		return errors.New("user already verified")
	}

	// Generate new OTP and update expiry
	user.Otp = helper.GenerateOTP()
	user.OtpExpiry = time.Now().Add(5 * time.Minute)
	user.UpdatedAt = time.Now()

	// Send the OTP via email and SMS
	if err := helper.SendOtpEmail(s.Cfg, user.Email, user.Otp); err != nil {
		log.Println("Failed to send OTP Email:", err)
	}
	if err := helper.SendOtpSMS(s.Cfg, user.Phone, user.Otp); err != nil {
		log.Println("Failed to send OTP SMS:", err)
	}

	// Update the user record
	if err := s.UserRepo.Update(ctx, user); err != nil {
		log.Println("Failed to update user with new OTP:", err)
		return err
	}

	log.Println("OTP sent successfully to :", email)
	return nil
}

func (s *AuthService) VerifyOtp(ctx context.Context, email, otp string) error {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}
	if user.IsVerified {
		return errors.New("user already verified")
	}

	if user.Otp != otp {
		return errors.New("invalid otp")
	}

	if time.Now().After(user.OtpExpiry) {
		return errors.New("OTP expired")
	}

	user.IsVerified = true
	user.Otp = ""
	user.OtpExpiry = time.Time{}
	if err := s.UserRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResendOtp(ctx context.Context, email string) error {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.IsVerified {
		return errors.New("user already verified")
	}

	// Prevnet too frequent resends
	if time.Until(user.OtpExpiry) > 4*time.Minute {
		return errors.New("please wait before requesting another OTP")
	}

	// Resend otp
	otp := helper.GenerateOTP()
	user.Otp = otp
	user.OtpExpiry = time.Now().Add(5 * time.Minute)

	if err = helper.SendOtpEmail(s.Cfg, user.Email, otp); err != nil {
		log.Println("Failed to send OTP Email:", err)
	}
	if err = helper.SendOtpSMS(s.Cfg, user.Phone, otp); err != nil {
		log.Println("Failed to send OTP SMS:", err)
	}

	log.Printf("Resent OTP to: %s", user.Email)
	return s.UserRepo.Update(ctx, user)
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Trim inputs for consistency
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	// Step 1: Fetch user by email
	user, err := a.UserRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		log.Println("User not found or DB error:", err)
		return "", errors.New("invalid credentials")
	}

	// Step 2: Check verification
	if !user.IsVerified {
		log.Println("Unverified user tried to login:", user.Email)
		return "", errors.New("please verify your email before login")
	}

	// Step 3: Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("Password mismatch for user:", user.Email)
		return "", errors.New("invalid credentials")
	}

	// Step 4: Generate JWT token
	token, err := jwtutil.GenerateToken(user.Email, "user", a.Cfg.JWTSecret, time.Hour*24)
	if err != nil {
		log.Println("Failed to generate token for user:", user.Email)
		return "", errors.New("internal server error")
	}

	log.Println("User logged in successfully:", user.Email)
	return token, nil
}
