package main

import (
	"fmt"
	"time"

	flowmailer "github.com/blutspende/go-flowmailer"
)

func main() {
	fm := flowmailer.New(6377, "870605337621aa15a8645cc3eb80e595b1c67d2a", "eec45fcc32a696e64c0f249a6c3c161dd1a8a80f")

	messages, maxitems, err := fm.GetMessagesHeld(time.Now().Add(-24*time.Hour), time.Now(), 0, 10)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, m := range messages {
		fmt.Printf("%+v\n", m)
	}

	fmt.Printf("%d messages total", maxitems)

}
