package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"
	
	// 生成哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", string(hash))
	
	// 验证哈希
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
	} else {
		fmt.Println("Verification successful!")
	}
	
	// 测试数据库中的哈希
	dbHash := "$2a$10$dOB8D9Opb.9fjm/AiYpJ8uAcBQGbz10AcTO.HQp3S3qYAwULXVfiS"
	fmt.Printf("\nTesting DB hash: %s\n", dbHash)
	err = bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
	if err != nil {
		fmt.Printf("DB hash verification failed: %v\n", err)
	} else {
		fmt.Println("DB hash verification successful!")
	}
}