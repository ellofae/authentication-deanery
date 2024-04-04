package utils

import (
	"crypto/rand"
)

func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	password := make([]byte, length)
	rand.Read(password)
	for i, b := range password {
		password[i] = charset[b%byte(len(charset))]
	}
	return string(password)
}
