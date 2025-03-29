package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hamidoujand/sales/cmd/admin/commands"
	"github.com/hamidoujand/sales/internal/sqldb"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("admin <subcommand> [...args]")
		fmt.Println("==========================SubCommands===========================")
		fmt.Println("genkey: generate a set of private/public key files.")
		fmt.Println("gentoken: generate a JWT token for userid using a key.")
		fmt.Println("migrate: migrates the database.")
		fmt.Println("=====================================================")
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("expected subcommands to be passed")
	}

	switch os.Args[1] {
	case "genkey":
		genkeyCommand := flag.NewFlagSet("genkey", flag.ExitOnError)
		keySize := genkeyCommand.Int("size", 2048, "key size in bits.")
		//parse the args
		genkeyCommand.Parse(os.Args[2:])
		if err := commands.GenerateKey(*keySize); err != nil {
			return fmt.Errorf("generateKey: %w", err)
		}
	case "gentoken":
		genTokenCommand := flag.NewFlagSet("gentoken", flag.ExitOnError)
		userID := genTokenCommand.String("userid", "", "id of the userbus that token will belong.")
		kid := genTokenCommand.String("kid", "", "ID of the private key used to sign the token.")
		keyPath := genTokenCommand.String("keypath", "infra/keys", "path to the dir the holds private and public key pairs.")

		genTokenCommand.Parse(os.Args[2:])

		if *userID == "" || *kid == "" {
			fmt.Println("Usage: gentoken kid=<key id> userid=<userbus id> [keypath=<path to keys folder>]")
			return errors.New("kid and userid are required")
		}

		if err := commands.GenerateToken(*keyPath, *userID, *kid); err != nil {
			fmt.Println("Usage: gentoken kid=<key id> userid=<userbus id> [keypath=<path to keys folder>]")
			return fmt.Errorf("generate token: %w", err)
		}
	case "migrate":
		migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
		user := migrateCommand.String("user", "postgres", "user is the database user.")
		pass := migrateCommand.String("pass", "password", "password for the database user.")
		host := migrateCommand.String("host", "localhost:5432", "database host.")
		db := migrateCommand.String("dbname", "postgres", "name of the db that you want to run migrations againt it.")

		migrateCommand.Parse(os.Args[2:])
		cfg := sqldb.Config{
			Host:       *host,
			Password:   *pass,
			User:       *user,
			Name:       *db,
			DisableTLS: true,
		}
		if err := commands.Migrate(cfg); err != nil {
			fmt.Println("Usage: migrate host=<db_host> user=<db_user> pass=<db_pass> dbname=<name_of_db>")
			return fmt.Errorf("migrate: %w", err)
		}

	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}

	return nil
}
