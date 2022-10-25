package flowmailer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type flowmailer struct {
	account_id               int
	client                   *resty.Client
	client_id, client_secret string
	token                    string
	token_valid_until        time.Time
}

func (fm *flowmailer) Login() error {

	fm.token = ""

	resp, err := fm.client.R().
		EnableTrace().
		SetFormData(map[string]string{
			"client_id":     fm.client_id,
			"client_secret": fm.client_secret,
			"grant_type":    "client_credentials",
			"scope":         "api",
		}).
		Post("https://login.flowmailer.net/oauth/token")

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case 200:
		var oAuthTokenResponse OAuthTokenResponse
		err := json.Unmarshal(resp.Body(), &oAuthTokenResponse)
		if err != nil {
			return err
		}
		fm.token = oAuthTokenResponse.Access_token
		fm.token_valid_until = time.Now().Add(time.Duration(oAuthTokenResponse.Expires_in) * time.Second)
	case 401:
		return fmt.Errorf("Unauthorized")
	default:
		return fmt.Errorf("unkown Status when logging in: %d", resp.StatusCode())
	}

	return nil
}

func (fm *flowmailer) GetMessages(from, until time.Time, rangemin, rangemax int) ([]Message, int, error) {
	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		SetHeader("Range", fmt.Sprintf("items=%d-%d", rangemin, rangemax)).
		Get(fmt.Sprintf("https://api.flowmailer.net/%d/messages;daterange=%s,%s?addevents=true&addtags=true",
			fm.account_id,
			from.Format("2006-01-02T15:04:05-0700"),
			until.Format("2006-01-02T15:04:05-0700")))

	if err != nil {
		return nil, 0, err
	}

	switch resp.StatusCode() {
	case 206: // Partial Content
		message := make([]Message, 0)
		err := json.Unmarshal(resp.Body(), &message)
		if err != nil {
			return nil, 0, err
		}
		maxPage := len(message)
		contentRange := resp.Header().Get("Content-Range")

		if strings.Contains(contentRange, "/") {
			maxPage, _ = strconv.Atoi(strings.Split(contentRange, "/")[1])
		}
		return message, maxPage, nil
	case 401:
		err := fm.Login()
		if err != nil {
			return nil, 0, err
		}
		return fm.GetMessages(from, until, rangemin, rangemax)
	default:
		return nil, 0, fmt.Errorf("unexpected return-code %d", resp.StatusCode())
	}
}

func (fm *flowmailer) GetMessagesBySource(sourceid int, from, until time.Time, rangemin, rangemax int) ([]Message, int, error) {

	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		SetHeader("Range", fmt.Sprintf("items=%d-%d", rangemin, rangemax)).
		Get(fmt.Sprintf("https://api.flowmailer.net/%d/sources/%d/messages;daterange=%s,%s?addevents=true&addtags=true",
			fm.account_id,
			sourceid,
			from.Format("2006-01-02T15:04:05-0700"),
			until.Format("2006-01-02T15:04:05-0700")))

	if err != nil {
		return nil, 0, err
	}

	switch resp.StatusCode() {
	case 206: // Partial Content
		message := make([]Message, 0)
		err := json.Unmarshal(resp.Body(), &message)
		if err != nil {
			return nil, 0, err
		}
		maxPage := len(message)
		contentRange := resp.Header().Get("Content-Range")

		if strings.Contains(contentRange, "/") {
			maxPage, _ = strconv.Atoi(strings.Split(contentRange, "/")[1])
		}
		return message, maxPage, nil
	case 401:
		err := fm.Login()
		if err != nil {
			return nil, 0, err
		}
		return fm.GetMessagesBySource(sourceid, from, until, rangemin, rangemax)
	default:
		return nil, 0, fmt.Errorf("unexpected return-code %d", resp.StatusCode())
	}
}

func (fm *flowmailer) GetMessagesHeld(from, until time.Time, rangemin, rangemax int) ([]MessageHold, int, error) {

	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		SetHeader("Range", fmt.Sprintf("items=%d-%d", rangemin, rangemax)).
		Get(fmt.Sprintf("https://api.flowmailer.net/%d/message_hold;daterange=%s,%s",
			fm.account_id,
			from.Format("2006-01-02T15:04:05-0700"),
			until.Format("2006-01-02T15:04:05-0700")))

	if err != nil {
		return nil, 0, err
	}

	switch resp.StatusCode() {
	case 206: // Partial Content
		message := make([]MessageHold, 0)
		err := json.Unmarshal(resp.Body(), &message)
		if err != nil {
			return nil, 0, err
		}
		maxPage := len(message)
		contentRange := resp.Header().Get("Content-Range")

		if strings.Contains(contentRange, "/") {
			maxPage, _ = strconv.Atoi(strings.Split(contentRange, "/")[1])
		}
		return message, maxPage, nil
	case 401:
		err := fm.Login()
		if err != nil {
			return nil, 0, err
		}
		return fm.GetMessagesHeld(from, until, rangemin, rangemax)
	default:
		return nil, 0, fmt.Errorf("unexpected return-code %d", resp.StatusCode())
	}
}

func (fm *flowmailer) SubmitEmail(toEmail, toName, fromEmail, fromName, subject, textBody, htmlBody string, attachments []Attachment) error {

	for i, _ := range attachments {
		if attachments[i].ContentType == "" {
			attachments[i].ContentType = "application/octet-stream"
		}
		if attachments[i].ContentId == "" { //default: a
			attachments[i].ContentId = uuid.New().String()
		}
		if attachments[i].Disposition == "" { // default: attachment
			attachments[i].Disposition = DISPOSITION_ATTACHMENT
		}
	}
	var sm SubmitMessage

	sm.HeaderFromAddress = fromEmail
	sm.HeaderFromName = fromName
	sm.HeaderToName = toName
	sm.MessageType = EMAIL
	sm.RecipientAddress = toEmail
	sm.SenderAddress = fromEmail
	sm.Subject = subject
	sm.Html = htmlBody
	sm.Text = textBody
	sm.Attachments = attachments

	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		SetBody(sm).
		Post(fmt.Sprintf("https://api.flowmailer.net/%d/messages/submit", fm.account_id))

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case 201:
		return nil
	case 401:
		err := fm.Login()
		if err != nil {
			return err
		}
		return fm.SubmitEmail(toEmail, toName, fromEmail, fromName, subject, textBody, htmlBody, attachments)
	default:
		return fmt.Errorf("unkown statuscode when sending email %d", resp.StatusCode())
	}

	return nil
}

func (fm *flowmailer) GetMessageFromArchiveById(id string) ([]MessageArchive, error) {
	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		Get(fmt.Sprintf("https://api.flowmailer.net/%d/messages/%s/archive?addattachments=true&adddata=true",
			fm.account_id, id))

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case 200:
		messages := make([]MessageArchive, 0)
		err := json.Unmarshal(resp.Body(), &messages)
		if err != nil {
			return nil, err
		}
		return messages, nil
	case 401:
		err := fm.Login()
		if err != nil {
			return nil, err
		}
		return fm.GetMessageFromArchiveById(id)
	default:
		return nil, fmt.Errorf("unexpected return-code %d", resp.StatusCode())
	}
}

func (fm *flowmailer) GetAttachmentFromArchiveMessage(messageId, flowStepId, contentId string) (*Attachment, error) {
	resp, err := fm.client.R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", fm.token)).
		SetHeader("Accept", "application/vnd.flowmailer.v1.12+json;charset=UTF-8").
		Get(fmt.Sprintf("https://api.flowmailer.net/%d/messages/%s/archive/%s/attachment/%s?addattachments=true&adddata=true",
			fm.account_id,
			messageId,
			flowStepId,
			contentId))

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case 200:
		var attachment Attachment
		err := json.Unmarshal(resp.Body(), &attachment)
		if err != nil {
			return nil, err
		}
		return &attachment, nil
	case 401:
		err := fm.Login()
		if err != nil {
			return nil, err
		}
		return fm.GetAttachmentFromArchiveMessage(messageId, flowStepId, contentId)
	default:
		return nil, fmt.Errorf("unexpected return-code %d", resp.StatusCode())
	}

}

func New(account_id int, client_id, client_secret string) Flowmailer {
	return &flowmailer{
		client:        resty.New(),
		client_id:     client_id,
		client_secret: client_secret,
		account_id:    account_id,
	}
}
