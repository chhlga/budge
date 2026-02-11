package email

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

// Parse parses raw email data into a Message struct
func Parse(data []byte) (*Message, error) {
	reader := bytes.NewReader(data)

	mr, err := mail.CreateReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail reader: %w", err)
	}

	msg := &Message{
		Body: &Body{},
	}

	// Parse headers
	header := mr.Header

	// From
	if addrs, err := header.AddressList("From"); err == nil && len(addrs) > 0 {
		msg.From = make([]Address, len(addrs))
		for i, addr := range addrs {
			msg.From[i] = Address{
				Name:  addr.Name,
				Email: addr.Address,
			}
		}
	}

	// To
	if addrs, err := header.AddressList("To"); err == nil && len(addrs) > 0 {
		msg.To = make([]Address, len(addrs))
		for i, addr := range addrs {
			msg.To[i] = Address{
				Name:  addr.Name,
				Email: addr.Address,
			}
		}
	}

	// Cc
	if addrs, err := header.AddressList("Cc"); err == nil && len(addrs) > 0 {
		msg.Cc = make([]Address, len(addrs))
		for i, addr := range addrs {
			msg.Cc[i] = Address{
				Name:  addr.Name,
				Email: addr.Address,
			}
		}
	}

	// Bcc
	if addrs, err := header.AddressList("Bcc"); err == nil && len(addrs) > 0 {
		msg.Bcc = make([]Address, len(addrs))
		for i, addr := range addrs {
			msg.Bcc[i] = Address{
				Name:  addr.Name,
				Email: addr.Address,
			}
		}
	}

	// Reply-To
	if addrs, err := header.AddressList("Reply-To"); err == nil && len(addrs) > 0 {
		msg.ReplyTo = make([]Address, len(addrs))
		for i, addr := range addrs {
			msg.ReplyTo[i] = Address{
				Name:  addr.Name,
				Email: addr.Address,
			}
		}
	}

	// Subject
	if subject, err := header.Subject(); err == nil {
		msg.Subject = subject
	}

	// Date
	if date, err := header.Date(); err == nil {
		msg.Date = date
	}

	// Message-ID
	if msgID, err := header.Text("Message-ID"); err == nil {
		msg.MessageID = msgID
	}

	// In-Reply-To
	if inReplyTo, err := header.Text("In-Reply-To"); err == nil {
		msg.InReplyTo = inReplyTo
	}

	// References
	if references, err := header.Text("References"); err == nil {
		// References can be space or newline separated
		refList := strings.Fields(references)
		msg.References = refList
	}

	// Content-Type
	contentType, _, _ := header.ContentType()
	msg.ContentType = contentType

	// Parse body parts
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read part: %w", err)
		}

		// Try to get content type from the part header
		contentTypeStr := part.Header.Get("Content-Type")
		partType := contentTypeStr
		if idx := strings.Index(contentTypeStr, ";"); idx > 0 {
			partType = contentTypeStr[:idx]
		}
		partType = strings.TrimSpace(partType)

		switch {
		case strings.HasPrefix(partType, "text/plain"):
			body, err := io.ReadAll(part.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read text/plain body: %w", err)
			}
			msg.Body.Text = string(body)

		case strings.HasPrefix(partType, "text/html"):
			body, err := io.ReadAll(part.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read text/html body: %w", err)
			}
			msg.Body.HTML = string(body)

		case strings.HasPrefix(partType, "multipart/"):
			// Multipart content will be handled by NextPart() iteration
			continue

		default:
			// Could be an attachment, but we'll handle that in future enhancements
			continue
		}
	}

	return msg, nil
}
