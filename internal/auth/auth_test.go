package auth

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	token, _ := HashPassword("admin")
	fmt.Println(token)
}
