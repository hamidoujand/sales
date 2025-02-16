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
		fmt.Println("=====================================================")
		fmt.Println("genkey: generate a set of private/public key files.")
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
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}

	return nil
}
