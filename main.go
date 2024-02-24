package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gopkg.in/yaml.v3"
)

type Email struct {
	Username string `yaml:"login"`
	Password string `yaml:"password"`
}

type Server struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

type Config struct {
	Server Server `yaml:"server"`
	Email  Email  `yaml:"email"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage %s <config file>", os.Args[0])
		os.Exit(1)
	}

	read, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read file %s (%v)", os.Args[1], err)
		os.Exit(1)
	}

	var config Config

	err = yaml.Unmarshal(read, &config)

	imap_name := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	client, err := client.DialTLS(imap_name, nil)
	if err != nil {
		log.Fatalf("could not connect to imap server: %s %e", imap_name, err)
	}
	defer client.Logout()

	fmt.Printf("%+v", config)

	if err := client.Login(config.Email.Username, config.Email.Password); err != nil {
		log.Fatalf("could not authenticate with Imap server %v", err)
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.List("", "*", mailboxes)

	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Printf("    - %s", m.Name)
	}

	if err := <-done; err != nil {
		log.Printf("error: %v\n", err)
	}

	_, err = client.Select("INBOX", true)
	if err != nil {
		log.Fatalf("Error selecint INBOX: %v", err)
	}

	searchCriteria := imap.NewSearchCriteria()
	searchCriteria.WithoutFlags = []string{imap.SeenFlag}
	seqnums, err := client.Search(searchCriteria)
	if err != nil {
		log.Printf("could not search message box: %v", err)
	}

	if len(seqnums) == 0 {
		fmt.Printf("0")
		os.Exit(0)
	}

	messages := make(chan *imap.Message)

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqnums...)

	done = make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{imap.FetchFull}, messages)

	}()

	log.Println("Messages:")
	for msg := range messages {
		log.Printf("    * %+v %s", msg.Flags, msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Printf("Error reading messages: %v", err)
	}

	fmt.Printf("%d", len(seqnums))
}
