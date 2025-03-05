package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/sec"
)

// right now the seed program will focus on seeding accounts but later we might add other resources
// and before each resource seed we ask if the user wants to seed the resource, if not, move to the next resource

type Account struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

type Court struct {
	Name string `json:"name"`
}

func main() {
	var seedAccountsFile string
	var seedCourtsFile string
	flag.StringVar(&seedAccountsFile, "accounts-file", "./db/seed/accounts.json", "Path to seed file")
	flag.StringVar(&seedCourtsFile, "courts-file", "./db/seed/courts.json", "Path to seed file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Seed accounts and courts into the database\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	ctx := context.Background()

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// connect db
	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var confirm string
	fmt.Printf("Do you want to seed accounts found at %s? (y/n): ", seedAccountsFile)
	fmt.Scanln(&confirm)
	if confirm == "y" {
		// read seed file
		b, err := os.ReadFile(seedAccountsFile)
		if err != nil {
			log.Fatal(err)
		}

		accounts := []Account{}
		err = json.Unmarshal(b, &accounts)
		if err != nil {
			log.Fatal(err)
		}

		err = seedAccounts(ctx, db, accounts)
		if err != nil {
			log.Fatalf("seeding accounts: %v", err)
		}

		fmt.Println("accounts seeded!")
	}

	fmt.Printf("Do you want to seed courts found at %s? (y/n): ", seedCourtsFile)
	fmt.Scanln(&confirm)
	if confirm == "y" {
		fmt.Println("seed courts")
	}
}

func seedAccounts(ctx context.Context, db *db.Conn, data []Account) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	accountSql := `
		insert into account (name, email, dob, gender, phone_number, password, role)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id
	`

	playerSql := `
		insert into player (account_id)
		values ($1)
	`

	for _, v := range data {
		enc, err := sec.EncryptPwd(v.Password)
		if err != nil {
			return err
		}
		v.Password = enc

		var accountID string
		err = tx.QueryRow(ctx, accountSql, v.Name, v.Email, v.Dob, v.Gender, v.PhoneNumber, v.Password, v.Role).Scan(&accountID)
		if err != nil {
			return fmt.Errorf("failed to insert account %s -> %w", v.Email, err)
		}

		_, err = tx.Exec(ctx, playerSql, accountID)
		if err != nil {
			return fmt.Errorf("failed to insert player for account %s -> %w", v.Email, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func seedCourts(ctx context.Context, db *db.Conn, data []Court) error {
	return nil
}
