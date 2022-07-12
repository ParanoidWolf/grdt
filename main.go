package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    /*
     * body := get_json("memes")
     * image_list := get_images(body)
     * fmt.Println(image_list)
    */
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
