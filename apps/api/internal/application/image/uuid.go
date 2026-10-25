package image

import "github.com/google/uuid"

func uuidV7() string {
	return uuid.Must(uuid.NewV7()).String()
}