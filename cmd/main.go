package main

import (
	"fmt"
	// "os"
	//
	// "dbtui/internal/render"
	"dbtui/internal/utils"

	// tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := utils.ParseArgs()

	fmt.Println(args.DBPath)
	// p := tea.NewProgram(
	// 	render.InitialModel(),
	// )
	//
	// if len(os.Getenv("DEBUG")) > 0 {
	// 	f, err := tea.LogToFile("debug.log", "debug")
	// 	if err != nil {
	// 		fmt.Println("fatal:", err)
	// 		os.Exit(1)
	// 	}
	// 	defer f.Close()
	// }
	//
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}
