package utils

import (
	"fmt"
	"log"
	"net/http"
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

func FormatDate(input string) string {
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", input) // Adjust the layout to match your input format
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return input // Return the original string if parsing fails
	}
	return parsedTime.Format("2006-01-02") // Format the date as YYYY-MM-DD
}

func ChainMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
