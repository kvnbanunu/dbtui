package utils

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Args struct {
	Help   bool
	Seed   bool
	DBPath string
}

func ParseArgs() *Args {
	args := Args{}
	flag.BoolVar(&args.Help, "h", false, "Displays this help message")
	flag.BoolVar(&args.Seed, "seed", false, "Seeds database with test data")
	flag.Parse()

	if args.Help {
		usage("")
	}

	remaining := flag.Args()
	if len(remaining) != 1 {
		usage("DB PATH is missing")
	}

	args.DBPath = remaining[0]

	return &args
}

func usage(msg string) {
	if msg != "" {
		log.Println(msg)
	}

	fmt.Printf(`Usage: dbtui [OPTIONS] <DB PATH>
Options:
	-h     Displays this help message
	-seed  Inserts dummy data into the database
`)
	os.Exit(1)
}
