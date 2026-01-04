package shared

import "github.com/google/uuid"

func GenerateId() string {
	newId, _ := uuid.NewV7()
	return newId.String()
}
