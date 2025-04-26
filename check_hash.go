package main

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	password := "123456aA@"

// 	// Generate hash with default cost (10)
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Printf("Error generating hash: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("New hash for password '%s': %s\n", password, string(hashedPassword))

// 	// Verify the hash
// 	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
// 	if err != nil {
// 		fmt.Printf("Verification failed: %v\n", err)
// 	} else {
// 		fmt.Println("Hash verified successfully!")
// 	}
// }
