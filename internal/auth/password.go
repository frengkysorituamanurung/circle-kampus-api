package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Angka 14 adalah "cost"
	return string(bytes), err
}

// CheckPasswordHash membandingkan password teks biasa dengan hash-nya
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}