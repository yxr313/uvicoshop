package image

import "github.com/google/uuid"

func mustUUID() string {
	return uuid.Must(uuid.NewV7()).String()
}