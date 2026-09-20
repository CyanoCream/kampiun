package domain

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID = id unik 32 hex.
func NewID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ShortCode = kode unik pendek (mis utk akses private viewer), base huruf besar + angka,
// tanpa karakter ambigu (O/0/I/1).
func ShortCode(n int) string {
	const alpha = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = alpha[int(b[i])%len(alpha)]
	}
	return string(b)
}
