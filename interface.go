package flowmailer

import "time"

type Flowmailer interface {
	// You do not need to call this directly. Login is called on expired sessions by each Method
	Login() error

	// Get all received messages
	GetMessages(from, until time.Time, rangemin, rangemax int) ([]Message, int, error)

	// Get a list of the held messages
	// Returns:
	//   - List of Messages
	//   - Count of Items in total (not only the pagination window)
	//   - nil if nor error occurs
	GetMessagesHeld(from, until time.Time, rangemin, rangemax int) ([]MessageHold, int, error)

	// Get a list of the messages by source
	// Returns:
	//   - List of Messages
	//   - Count of Items in total (not only the pagination window)
	//   - nil if nor error occurs
	GetMessagesBySource(sourceid int, from, until time.Time, rangemin, rangemax int) ([]Message, int, error)

	// Get a message that was archived by a flow (required if you want the attachment)
	GetMessageFromArchiveById(id string) ([]MessageArchive, error)

	// Get the attachment of a message retrieved with GetMessageFromArchiveById before
	GetAttachmentFromArchiveMessage(messageId, flowStepId, contentId string) (*Attachment, error)

	// Send a new Email
	SubmitEmail(toEmail, toName, fromEmail, fromName, subject, textbody, htmlbody string, attachments []Attachment) error
}
