package main

import (
	"fmt"
	"log"
	"os"

	"dbtui/internal/render"
	"dbtui/internal/database"
	"dbtui/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := utils.ParseArgs()

	db, err := database.Init(args.DBPath)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	defer db.Close()

	err = db.InsertDummy()
	if err != nil {
		log.Fatalln("Error inserting dummy data:", err)
	}

	initialModel, err := render.InitialModel(db)
	if err != nil {
		log.Fatalln("Error getting initial model:", err)
	}

	p := tea.NewProgram(initialModel)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
