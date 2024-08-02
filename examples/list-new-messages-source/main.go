package main

import (
	"fmt"
	"time"

	flowmailer "github.com/blutspende/go-flowmailer"
)

func main() {
	fm := flowmailer.New(6377, "870605337621aa15a8645cc3eb80e595b1c67d2a", "eec45fcc32a696e64c0f249a6c3c161dd1a8a80f")

	err := fm.Login()
	if err != nil {
		fmt.Println(err.Error())
	}

	messages, maxitems, err := fm.GetMessagesBySource(26177, time.Now().Add(-24*time.Hour), time.Now(), 0, 1)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%#v\n", messages)
	fmt.Printf("%d messages in total available\n", maxitems)
}
