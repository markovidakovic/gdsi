// seed is used one time for seeding the database with a developer account
// format of arguments arg1=val1 arg2=val2
package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/security"
	"github.com/markovidakovic/gdsi/server/validation"
)

func main() {
	// get cli arguments
	var args []string = os.Args[1:]
	if len(args) == 0 {
		log.Fatal("no arguments provided")
	}

	var name, email, dob, gender, phone, password string
	var role string = "developer"

	// process cli arguments
	for _, arg := range args {
		argSl := strings.Split(arg, "=")

		switch argSl[0] {
		case "name":
			name = argSl[1]
		case "email":
			email = argSl[1]
		case "dob":
			dob = argSl[1]
		case "gender":
			gender = argSl[1]
		case "phone":
			phone = argSl[1]
		case "password":
			password = argSl[1]
		default:
			log.Fatalf("argument %q not recognized", argSl[1])
		}
	}

	// validate input
	if name == "" {
		log.Fatal("argument name empty")
	}
	if email == "" {
		log.Fatalf("argument email empty")
	} else if !validation.IsValidEmail(email) {
		log.Fatalf("argument email is invalid")
	}
	if dob == "" {
		log.Fatalf("argument dob empty")
	} else {
		if _, err := time.Parse("2006-01-02", dob); err != nil {
			log.Fatal("argument dob invalid")
		}
	}
	if gender == "" {
		log.Fatal("argument gender empty")
	} else if gender != "male" && gender != "female" {
		log.Fatal("argument gender invalid")
	}
	if phone == "" {
		log.Fatal("argument phone empty")
	} else if !validation.IsValidPhone(phone) {
		log.Fatal("argument phone invalid")
	}
	if password == "" {
		log.Fatal("argument password empty")
	}

	password, err := security.EncryptPwd(password)
	if err != nil {
		log.Fatal("failed to encrypt password")
	}

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config")
	}

	ctx := context.Background()

	// connect to the db
	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to the database")
	}

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

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
		log.Println("developer account seeded")
	}()

	err = db.QueryRow(ctx, sql1, name, email, dob, gender, phone, password, role).Scan(&dest)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(ctx, sql2, dest)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
