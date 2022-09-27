package hsh

import (
	"crypto/sha1"
	"fmt"
	"log"
)

type HashService struct {
}

func GetHashService() *HashService {
	return &HashService{}
}

func (h *HashService) GenerateHashPassword(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	log.Println("Hash password was generated")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
