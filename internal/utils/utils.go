package utils

import "flag"

type Args struct {
	DBPath string
}

func ParseArgs() *Args {
	args := Args{}
	flag.StringVar(&args.DBPath, "p", "./database.sqlite", "Path to your SQLite database file")
	flag.Parse()

	return &args
}
