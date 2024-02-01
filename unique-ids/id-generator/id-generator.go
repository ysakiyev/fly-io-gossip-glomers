package id_generator

import (
	"github.com/google/uuid"
)

func GenerateId() string {
	return uuid.New().String()
}
