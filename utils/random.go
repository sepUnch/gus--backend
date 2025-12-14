package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString membuat string acak dengan panjang n
// Karakter terdiri dari Huruf Kapital dan Angka
func GenerateRandomString(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	
	// Seed generator angka acak agar hasilnya selalu beda tiap detik
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}