package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/rand"
)

func GenerateUUID() string {
	return uuid.NewString()
}

func GenerateULIDWithDash() (string, error) {
	// Create a source of randomness
	t := time.Now()
	entropy := rand.New(rand.NewSource(uint64(t.UnixNano())))

	// Generate the ULID
	newULID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// Validate the ULID length
	if len(newULID) < 26 {
		fmt.Print("Error to Generate string")
		GenerateULIDWithDash()
	}

	// Insere os hífens nas posições desejadas
	formatted := newULID[:8] + "-" + newULID[8:17] + "-" + newULID[17:]
	return formatted, nil
}
