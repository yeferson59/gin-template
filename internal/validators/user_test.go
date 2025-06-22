package validators

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Valid username", "testuser", false},
		{"Valid with underscore", "test_user", false},
		{"Valid with hyphen", "test-user", false},
		{"Valid with numbers", "user123", false},
		{"Empty username", "", true},
		{"Too short", "ab", true},
		{"Too long", "averylongusernamethatisgreaterthan30characters", true},
		{"With spaces", "test user", true},
		{"With special chars", "test@user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid with subdomain", "test@mail.example.com", false},
		{"Valid with plus", "test+tag@example.com", false},
		{"Valid with dash", "test-user@example.com", false},
		{"Empty email", "", true},
		{"Missing @", "testexample.com", true},
		{"Missing domain", "test@", true},
		{"Missing local part", "@example.com", true},
		{"Invalid format", "test@@example.com", true},
		{"Too long", func() string {
			// Create an email longer than 254 characters
			// "@example.com" = 12 characters, so we need 243+ for local part
			longLocal := make([]byte, 243)
			for i := range longLocal {
				longLocal[i] = 'a'
			}
			return string(longLocal) + "@example.com" // This will be 255 characters
		}(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "Password123!", false},
		{"Valid complex", "MySecure@Pass1", false},
		{"Empty password", "", true},
		{"Too short", "Pass1!", true},
		{"No uppercase", "password123!", true},
		{"No lowercase", "PASSWORD123!", true},
		{"No numbers", "Password!", true},
		{"No special chars", "Password123", true},
		{"Too long", func() string {
			// Create a password longer than 128 characters
			longPass := make([]byte, 130)
			for i := range longPass {
				longPass[i] = 'A'
			}
			longPass[129] = '1' // Add a number
			longPass[128] = '!' // Add a special character
			return string(longPass)
		}(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserRegistration(t *testing.T) {
	tests := []struct {
		name    string
		req     *AuthRequest
		wantErr bool
	}{
		{
			"Valid registration",
			&AuthRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			false,
		},
		{
			"Invalid username",
			&AuthRequest{
				Username: "ab",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			true,
		},
		{
			"Invalid email",
			&AuthRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "Password123!",
			},
			true,
		},
		{
			"Invalid password",
			&AuthRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "weak",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserRegistration(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserRegistration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserLogin(t *testing.T) {
	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
	}{
		{
			"Valid login",
			&LoginRequest{
				Username: "testuser",
				Password: "anypassword",
			},
			false,
		},
		{
			"Invalid username",
			&LoginRequest{
				Username: "ab",
				Password: "anypassword",
			},
			true,
		},
		{
			"Empty password",
			&LoginRequest{
				Username: "testuser",
				Password: "",
			},
			true,
		},
		{
			"Whitespace password",
			&LoginRequest{
				Username: "testuser",
				Password: "   ",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserLogin(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
