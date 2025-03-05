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
	var accountsFile, courtsFile string
	flag.StringVar(&accountsFile, "accounts-file", "./db/seed/accounts.json", "Path to seed file")
	flag.StringVar(&courtsFile, "courts-file", "./db/seed/courts.json", "Path to seed file")
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
	fmt.Printf("Do you want to seed accounts found at %s? (y/n): ", accountsFile)
	fmt.Scanln(&confirm)
	if confirm == "y" {
		b, err := os.ReadFile(accountsFile)
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

	fmt.Printf("Do you want to seed courts found at %s? (y/n): ", courtsFile)
	fmt.Scanln(&confirm)
	if confirm == "y" {
		b, err := os.ReadFile(courtsFile)
		if err != nil {
			log.Fatal(err)
		}

		courts := []Court{}
		err = json.Unmarshal(b, &courts)
		if err != nil {
			log.Fatal(err)
		}

		err = seedCourts(ctx, db, courts)
		if err != nil {
			log.Fatalf("seeding courts: %v", err)
		}

		fmt.Println("courts seeded!")
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
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	// defer func() {
	// 	err := tx.Rollback(ctx)
	// 	if err != nil && err != pgx.ErrTxClosed {
	// 		log.Printf("rollback failed: %v", err)
	// 	}
	// }()

	sql := `
		insert into court (name)
		values ($1)
	`

	for _, v := range data {
		_, err := tx.Exec(ctx, sql, v.Name)
		if err != nil {
			return fmt.Errorf("failed to insert court %s -> %w", v.Name, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
