package main

import (
	"bufio"
	"estiam/dictionary"
	"fmt"
	"os"
	"strings"
)

func main() {
	d, err := dictionary.New("dictionary.json")
	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return
	}
	done := make(chan error)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command (add/define/remove/list/exit): ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		command = strings.TrimSpace(command)

		switch command {
		case "add":
			go func() {
				actionAdd(d, reader, done)
			}()
			if err := <-done; err != nil {
				fmt.Println("Error adding entry:", err)
			}

		case "define":
			actionDefine(d, reader)
		case "remove":
			go func() {
				actionRemove(d, reader, done)
			}()
			if err := <-done; err != nil {
				fmt.Println("Error removing entry:", err)
			}
		case "list":
			actionList(d)
		case "exit":
			fmt.Println("Exiting program.")
			os.Exit(0)
		default:
			fmt.Println("Invalid command. Try again.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader, done chan<- error) {
	fmt.Print("Enter word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition, done)
	fmt.Println("Word added successfully.")
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to define: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Def: %s\n", entry)
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader, done chan<- error) {
	fmt.Print("Enter word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word, done)
	fmt.Println("Word removed successfully.")
}

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()

	fmt.Println("Words:")
	for _, word := range words {
		fmt.Println(word)
	}

	fmt.Println("\nDictionary:")
	for word, entry := range entries {
		fmt.Printf("%s: %s\n", word, entry)
	}
}
