package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hamidoujand/sales/cmd/admin/commands"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("admin <subcommand> [...args]")
		fmt.Println("==========================SubCommands===========================")
		fmt.Println("genkey: generate a set of private/public key files.")
		fmt.Println("gentoken: generate a JWT token for userid using a key.")
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
		userID := genTokenCommand.String("userid", "", "id of the user that token will belong.")
		kid := genTokenCommand.String("kid", "", "ID of the private key used to sign the token.")
		keyPath := genTokenCommand.String("keypath", "infra/keys", "path to the dir the holds private and public key pairs.")

		genTokenCommand.Parse(os.Args[2:])

		if *userID == "" || *kid == "" {
			fmt.Println("Usage: gentoken kid=<key id> userid=<user id> [keypath=<path to keys folder>]")
			return errors.New("kid and userid are required")
		}

		if err := commands.GenerateToken(*keyPath, *userID, *kid); err != nil {
			fmt.Println("Usage: gentoken kid=<key id> userid=<user id> [keypath=<path to keys folder>]")
			return fmt.Errorf("generate token: %w", err)
		}

	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}

	return nil
}
