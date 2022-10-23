package main

import (
	"encoding/json"
	"fmt"
	"time"

	flowmailer "github.com/DRK-Blutspende-BaWueHe/go-flowmailer"
)

func main() {
	fm := flowmailer.New(6377, "870605337621aa15a8645cc3eb80e595b1c67d2a", "eec45fcc32a696e64c0f249a6c3c161dd1a8a80f")

	err := fm.Login()
	if err != nil {
		fmt.Println(err.Error())
	}

	messages, _, err := fm.GetMessages(time.Now().Add(-24*time.Hour), time.Now(), 0, 10)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Using message %s - %s\n", messages[0].SenderAddress, messages[0].Subject)

	msg, err := fm.GetMessageFromArchiveById(messages[0].Id)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%#v\n", msg)

	for _, at := range msg[0].Attachments {
		atdata, err := fm.GetAttachmentFromArchiveMessage(messages[0].Id, msg[0].FlowStepId, at.ContentId)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Attachmentdata: %#v\n", atdata)
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
