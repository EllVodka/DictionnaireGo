package main

import (
	"flag"
	"fmt"
	"os"

	"training.go/Dictionary/dictionary"
)

func main() {
	action := flag.String("action", "list", "action a faire sur le dictionaire")

	d, err := dictionary.New("./badger")
	handleErr(err)
	defer d.Close()

	flag.Parse()
	switch *action {
	case "list":
		actionList(d)
	case "add":
		actionAdd(d, flag.Args())
	case "define":
		actionDefine(d, flag.Arg(0))
	case "remove":
		actionRemove(d, flag.Arg(0))
	default:
		fmt.Printf("action inconue: %v\n", action)
	}
}

func actionDefine(d *dictionary.Dictionary, word string) {
	entries, err := d.Get(word)
	handleErr(err)
	fmt.Println(entries)
}

func actionRemove(d *dictionary.Dictionary, word string) {
	err := d.Remove(word)
	handleErr(err)
	fmt.Printf("'%v' est bien retiré du dictionaire", word)
}

func actionList(d *dictionary.Dictionary) {
	words, entries, err := d.List()
	handleErr(err)
	fmt.Println("Contenu du dictionaire")
	for _, word := range words {
		fmt.Println(entries[word])
	}
}

func actionAdd(d *dictionary.Dictionary, args []string) {
	word := args[0]
	definition := args[1]
	err := d.Add(word, definition)
	handleErr(err)
	fmt.Printf("'%v' ajouter au dictionaire\n", word)

}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Erreur de dictinaire: %v\n", err)
		os.Exit(1)
	}

}
