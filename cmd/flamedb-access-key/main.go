package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/myzie/flamedb/database"
	"github.com/namsral/flag"
)

func fatal(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
	os.Exit(1)
}

func main() {

	var (
		name  string
		refID string
		perm  string
	)

	flag.StringVar(&name, "name", "", "Key name")
	flag.StringVar(&refID, "ref", "", "Reference ID")
	flag.StringVar(&perm, "perm", "rw", "Permission: r or rw")
	flag.Parse()

	if name == "" {
		fatal("Please provide a key name using -name", nil)
	}

	gormDB, err := database.Connect(database.GetSettings())
	if err != nil {
		fatal("Failed to connect to database", err)
	}
	if err := gormDB.AutoMigrate(&database.AccessKey{}).Error; err != nil {
		fatal("Failed to migrate database", err)
	}

	permission := database.AccessKeyPermission(perm)
	key, plainText, err := database.NewAccessKey(name, refID, permission)
	if err != nil {
		fatal("Failed to generate access key", err)
	}
	if err := gormDB.Save(&key).Error; err != nil {
		fatal("Failed to save access key", err)
	}

	var output = struct {
		KeyID     string `json:"key_id"`
		KeySecret string `json:"key_secret"`
	}{
		KeyID:     key.ID,
		KeySecret: plainText,
	}

	js, err := json.Marshal(output)
	if err != nil {
		fatal("Failed to marshal output", err)
	}

	fmt.Printf("%s\n", string(js))
}
