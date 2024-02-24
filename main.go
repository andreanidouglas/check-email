package main

import (
	"fmt"

	imap "github.com/emersion/go-imap"
)

func main() {

	fmt.Println("Hello World")

	imap_name  := "imap.google.com" 

	client, err := imap.DialTls()
}
