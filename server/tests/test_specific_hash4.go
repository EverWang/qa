package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 从数据库中获取的具体哈希
	dbHash := "$2a$10$/Xa0ndJQ.pKk2095hL3AHeDGRlcQK1eyk4xHjq7fqOKIRuIEtRXHy"
	password := "123456"
	
	fmt.Printf("Testing specific hash from database:\n")
	fmt.Printf("Hash: %s\n", dbHash)
	fmt.Printf("Password: %s\n", password)
	
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
	if err != nil {
		fmt.Printf("Verification FAILED: %v\n", err)
	} else {
		fmt.Printf("Verification SUCCESSFUL!\n")
	}
	
	// 生成新的哈希进行对比
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating new hash: %v\n", err)
		return
	}
	
	fmt.Printf("\nNew hash for comparison: %s\n", string(newHash))
	
	// 验证新生成的哈希
	err = bcrypt.CompareHashAndPassword(newHash, []byte(password))
	if err != nil {
		fmt.Printf("New hash verification FAILED: %v\n", err)
	} else {
		fmt.Printf("New hash verification SUCCESSFUL!\n")
	}
}