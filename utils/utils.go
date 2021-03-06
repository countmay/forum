package utils

import (
	"fmt"
	"math/rand"
)

// Generate ...
func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
