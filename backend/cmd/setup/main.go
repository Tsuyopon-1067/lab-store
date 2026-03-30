package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"purchase-system/db"
	"purchase-system/service"
)

func main() {
	// DB初期化
	database, err := db.Init("../data/purchase.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// パスワード入力
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter admin password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("Error: Password cannot be empty")
		return
	}

	// パスワード確認
	fmt.Print("Confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)

	if password != confirmPassword {
		fmt.Println("Error: Passwords do not match")
		return
	}

	// パスワード設定
	authService := service.NewAuthService(database)
	if err := authService.SetPassword(password); err != nil {
		log.Fatalf("Failed to set password: %v", err)
	}

	fmt.Println("Admin password has been set successfully!")
}
