// seed is used one time for seeding the database with a developer account
// format of arguments -arg1="val1" -arg2="val2"
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/sec"
)

type input struct {
	name     string
	email    string
	dob      string
	gender   string
	phone    string
	password string
	role     string
}

func main() {
	// parse cli args
	input := parseInput()

	// validate input
	err := validateInput(input)
	if err != nil {
		log.Fatalf("invalid input: %v", err)
	}

	// encrypt pwd
	input.password, err = sec.EncryptPwd(input.password)
	if err != nil {
		log.Fatalf("encrypting password: %v", err)
	}

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	// connect to the db
	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("connecting database: %v", err)
	}

	ctx := context.Background()

	// seed dev account
	err = seedDeveloperAccount(ctx, db, input)
	if err != nil {
		log.Fatalf("seeding developer account: %v", err)
	}

}

func parseInput() input {
	in := input{}
	in.role = "developer"

	flag.StringVar(&in.name, "name", "", "Full name")
	flag.StringVar(&in.email, "email", "", "Email")
	flag.StringVar(&in.dob, "dob", "", "Date of birth")
	flag.StringVar(&in.gender, "gender", "", "Gender (male or female)")
	flag.StringVar(&in.phone, "phone", "", "Phone number")
	flag.StringVar(&in.password, "password", "", "Password")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Seed developer account into the database\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	return in
}

func validateInput(in input) error {
	if in.name == "" {
		return fmt.Errorf("name is required")
	}
	if in.email == "" {
		return fmt.Errorf("email is required")
	} else if !sec.IsValidEmail(in.email) {
		return fmt.Errorf("email is invalid")
	}
	if in.dob == "" {
		return fmt.Errorf("dob is required")
	} else {
		if _, err := time.Parse("2006-01-02", in.dob); err != nil {
			return fmt.Errorf("dob is invalid")
		}
	}
	if in.gender == "" {
		return fmt.Errorf("gender is required")
	} else if in.gender != "male" && in.gender != "female" {
		return fmt.Errorf("gender is invalid, expected male or female")
	}
	if in.phone == "" {
		return fmt.Errorf("phone number is required")
	} else if !sec.IsValidPhone(in.phone) {
		return fmt.Errorf("phone number is invalid")
	}
	if in.password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func seedDeveloperAccount(ctx context.Context, db *db.Conn, in input) error {
	// queries
	sql1 := `
		insert into account (name, email, dob, gender, phone_number, password, role)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id
	`
	sql2 := `
		insert into player (account_id)
		values ($1)
	`

	var dest string

	// begin tx
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
		log.Println("developer account seeded")
	}()

	// account
	err = db.QueryRow(ctx, sql1, in.name, in.email, in.dob, in.gender, in.phone, in.password, in.role).Scan(&dest)
	if err != nil {
		return err
	}

	// player
	_, err = db.Exec(ctx, sql2, dest)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
