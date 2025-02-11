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

func (i *input) parse() {
	flag.StringVar(&i.name, "name", "Marko Vidaković", "Full name")
	flag.StringVar(&i.email, "email", "marko.vidakovic@gmail.com", "Email")
	flag.StringVar(&i.dob, "dob", "1995-09-04", "Date of birth")
	flag.StringVar(&i.gender, "gender", "male", "Gender (male or female)")
	flag.StringVar(&i.phone, "phone", "0989559516", "Phone number")
	flag.StringVar(&i.password, "password", "string", "Password")
	i.role = "developer"

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Seed developer account into the database\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func (i input) validate() error {
	if i.name == "" {
		return fmt.Errorf("name is required")
	}
	if i.email == "" {
		return fmt.Errorf("email is required")
	} else if !sec.IsValidEmail(i.email) {
		return fmt.Errorf("email is invalid")
	}
	if i.dob == "" {
		return fmt.Errorf("dob is required")
	} else {
		if _, err := time.Parse("2006-01-02", i.dob); err != nil {
			return fmt.Errorf("dob is invalid")
		}
	}
	if i.gender == "" {
		return fmt.Errorf("gender is required")
	} else if i.gender != "male" && i.gender != "female" {
		return fmt.Errorf("gender is invalid, expected male or female")
	}
	if i.phone == "" {
		return fmt.Errorf("phone number is required")
	} else if !sec.IsValidPhone(i.phone) {
		return fmt.Errorf("phone number is invalid")
	}
	if i.password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func main() {
	inp := input{}

	// parse cli args
	inp.parse()

	err := inp.validate()
	if err != nil {
		log.Fatalf("invalid input: %v", err)
	}

	// confirmation prompt
	fmt.Printf("about to seed developer account with email %s. continue? (y/n): ", inp.email)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		log.Fatal("operation aborted")
	}

	// encrypt pwd
	inp.password, err = sec.EncryptPwd(inp.password)
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
	err = seedDeveloperAccount(ctx, db, inp)
	if err != nil {
		log.Fatalf("seeding developer account: %v", err)
	}
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
		if tx != nil {
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
