package main

import (
	b64 "encoding/base64"
	"fmt"

	flowmailer "github.com/DRK-Blutspende-BaWueHe/go-flowmailer"
)

func main() {
	fm := flowmailer.New(6377, "870605337621aa15a8645cc3eb80e595b1c67d2a", "eec45fcc32a696e64c0f249a6c3c161dd1a8a80f")

	attachments := make([]flowmailer.Attachment, 0)
	var attachment flowmailer.Attachment
	attachment.Filename = "etwas.txt"
	attachment.Content = b64.StdEncoding.EncodeToString([]byte("Guten Tag"))
	// attachment.ContentType = "application/octet-stream"
	// attachment.ContentId = "abcdefg"
	// attachment.Disposition = DISPOSITION_ATTACHMENT
	attachments = append(attachments, attachment)

	err := fm.SubmitEmail("kuhr@posteo.de", "stephan Kuhr", "test@bsd-screeninglabor.de", "Stephan Kuhr", "Subject is this", "textbody", "htmlbody", attachments)
	if err != nil {
		fmt.Println(err.Error())
	}

}
