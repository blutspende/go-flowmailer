package flowmailer

import (
	"time"
)

type Object struct {
	// not implemented
}

type Disposition string

const DISPOSITION_ATTACHMENT Disposition = "attachment"
const DISPOSITION_INLINE Disposition = "inline"
const DISPOSITION_RELATED Disposition = "related"

type Attachment struct {
	Content     string      `json:"content,omitempty"`
	ContentId   string      `json:"contentId,omitempty"`   // Content-ID header (required for disposition related) Example: <part1.DE1D8F7E.E51807FF@flowmailer.com>
	ContentType string      `json:"contentType,omitempty"` // Examples: application/pdf, image/jpeg
	Disposition Disposition `json:"disposition,omitempty"` // Supported values include: attachment, inline and related, special value related should be used for images referenced in the HTML
	Filename    string      `json:"filename,omitempty"`
}
type SubmitMessage struct {
	Attachments              []Attachment `json:"attachments,omitempty"`
	Data                     Object       `json:"data,omitempty"`
	DeliveryNotificationType string       `json:"deliveryNotificationType,omitempty"`
	FlowSelector             string       `json:"flowSelector,omitempty"`
	HeaderFromAddress        string       `json:"headerFromAddress,omitempty"`
	HeaderFromName           string       `json:"headerFromName,omitempty"`
	HeaderToAddress          string       `json:"headerToAddress,omitempty"`
	HeaderToName             string       `json:"headerToName,omitempty"`
	Headers                  []Header     `json:"headers,omitempty"`
	Html                     string       `json:"html,omitempty"`
	MessageType              MessageType  `json:"messageType,omitempty"`
	Mimedata                 string       `json:"mimedata,omitempty"`
	RecipientAddress         string       `json:"recipientAddress,omitempty"`
	ScheduleAt               time.Time    `json:"scheduleAt,omitempty"`
	SenderAddress            string       `json:"senderAddress,omitempty"`
	Subject                  string       `json:"subject,omitempty"`
	Tags                     []string     `json:"tags,omitempty"`
	Text                     string       `json:"text,omitempty"`
}

type MessageEvent struct {
	Data           string `json:"data,omitempty"`           // base64 Event data
	DeviceCategory string `json:"deviceCategory,omitempty"` //
	//ExtraData              string    `json:"extraData"`      // Event data
	Id                     string    `json:"id,omitempty"`       //Message event ID
	Inserted               time.Time `json:"inserted,omitempty"` // Database insert date
	LinkName               string    `json:"linkName,omitempty"`
	LinkTarget             string    `json:"linkTarget,omitempty"`
	MessageId              string    `json:"messageId,omitempty"`   // Message ID
	MessageTags            []string  `json:"messageTags,omitempty"` // Message tags- Only filled for the GET /{account_id}/message_events api call when the parameter addmessagetags is true
	Mta                    string    `json:"mta,omitempty"`         // mTA that reported this event
	OperatingSystem        string    `json:"operatingSystem,omitempty"`
	OperatingSystemVersion string    `json:"operatingSystemVersion,omitempty"`
	Received               time.Time `json:"received,omitempty"` // Event date
	Referer                string    `json:"referer,omitempty"`
	RemoteAddr             string    `json:"remoteAddr,omitempty"`
	Snippet                string    `json:"snippet,omitempty"`   // Bounce snippet or SMTP conversation snippet
	SubType                string    `json:"subType,omitempty"`   // Bounce sub type
	Tag                    string    `json:"tag,omitempty"`       // Custom event type
	EventType              string    `json:"eventType,omitempty"` // Event type, must be CUSTOM
	UserAgent              string    `json:"userAgent,omitempty"`
	UserAgentDisplayName   string    `json:"userAgentDisplayName,omitempty"`
	UserAgentString        string    `json:"userAgentString,omitempty"`
	UserAgentType          string    `json:"userAgentType,omitempty"`
	UserAgentVersion       string    `json:"userAgentVersion,omitempty"`
}

type ObjectDescription struct {
	Description string `json:"description,omitempty"`
	Id          string `json:"id,omitempty"`
}

type Address struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type MessageType string

const (
	EMAIL  MessageType = "EMAIL"
	SMS    MessageType = "SMS"
	LETTER MessageType = "LETTER"
)

type Message struct {
	BackendDone        time.Time         `json:"backendDone,omitempty"`        // The time flowmailer was done processing this message
	BackendStart       time.Time         `json:"backendStart,omitempty"`       // The time flowmailer started processing this message
	Events             []MessageEvent    `json:"events,omitempty"`             //Message events Ordered by received, new events first.
	Flow               ObjectDescription `json:"flow,omitempty"`               // Flow this message was processed in
	From               string            `json:"from,omitempty"`               // The email address in From email header
	FromAddress        Address           `json:"fromAddress,omitempty"`        // The address in From email header
	HeadersIn          []Header          `json:"headersIn,omitempty"`          // E-Mail headers of the submitted email message., Only applicable when messageType = EMAIL and addheaders parameter is true
	HeadersOut         []Header          `json:"headersOut,omitempty"`         // Headers of the final e-mail. Only applicable when messageType = EMAIL and addheaders parameter is true
	Id                 string            `json:"id,omitempty"`                 // Message id
	MessageDetailsLink string            `json:"messageDetailsLink,omitempty"` // Link for the message details page. With resend button.
	MessageIdHeader    string            `json:"messageIdHeader,omitempty"`    // Content of the Message-ID email header
	MessageType        MessageType       `json:"messageType,omitempty"`        // Message type: EMAIL, SMS or LETTER
	OnlineLink         string            `json:"onlineLink,omitempty"`         // Last online link
	RecipientAddress   string            `json:"recipientAddress,omitempty"`   // Recipient address
	SenderAddress      string            `json:"senderAddress,omitempty"`      // Sender address
	Source             ObjectDescription `json:"source,omitempty"`             // Source system that submitted this message
	Status             string            `json:"status,omitempty"`             // Current message status
	Subject            string            `json:"subject,omitempty"`            // Message subject Only applicable when messageType = EMAIL
	Submitted          time.Time         `json:"submitted,omitempty"`          // The time this message was submitted to flowmailer
	Tags               []string          `json:"tags,omitempty"`               // Message tags, only available for api calls with addtags=true
	ToAddressList      []Address         `json:"toAdressList,omitempty"`       // The recipients in the To email header
	TransactionId      string            `json:"transactionId,omitempty"`      // The SMTP transaction id, returned with the SMTP 250 response
}

type OAuthTokenResponse struct {
	Access_token string `json:"access_token,omitempty"`
	Expires_in   int    `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Token_type   string `json:"token_type,omitempty"`
}

type MessageHold struct {
	BackendDone   time.Time         `json:"backendDone,omitempty"`   // The time flowmailer was done processing this message
	Data          string            `json:"data,omitempty"`          // MIME message data or text for SMS messages
	DataCoding    byte              `json:"dataCoding,omitempty"`    // Only for SMS messages
	ErrorText     string            `json:"errorText,omitempty"`     // Message error text
	Flow          ObjectDescription `json:"flow,omitempty"`          // Flow this message was processed in
	MessageId     string            `json:"messageId,omitempty"`     // Message id
	MessageType   MessageType       `json:"messageType,omitempty"`   // Message type: EMAIL, SMS or LETTER
	Reason        string            `json:"reason,omitempty"`        // Message processing failure reason
	Recipient     string            `json:"recipient,omitempty"`     // Recipient address
	Sender        string            `json:"sender,omitempty"`        // The email address in From email header
	Source        ObjectDescription `json:"source,omitempty"`        // Source system that submitted this message
	Status        string            `json:"status,omitempty"`        // Current message status
	Submitted     time.Time         `json:"submitted,omitempty"`     // The time this message was submitted to flowmailer
	TransactionId string            `json:"transactionId,omitempty"` // The SMTP transaction id, returned with the SMTP 250 response
}

type MessageArchive struct {
	Attachments        []Attachment      `json:"attachments,omitempty"`
	Data               map[string]string `json:"data,omitempty"`
	FlowStepId         string            `json:"flowStepId,omitempty"`
	Html               string            `json:"html,omitempty"`
	MessageDetailsLink string            `json:"messageDetailsLink,omitempty"`
	MessageType        MessageType       `json:"messageType,omitempty"` // Message type: EMAIL, SMS or LETTER
	OnlineLink         string            `json:"onlineLink,omitempty"`
	OnlineVersion      bool              `json:"onlineVersion,omitempty"`
	Subject            string            `json:"subject,omitempty"`
	Text               string            `json:"text,omitempty"`
}
