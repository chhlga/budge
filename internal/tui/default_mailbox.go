package tui

import "strings"

func pickDefaultMailbox(mailboxes []string, configuredDefault string) string {
	configured := strings.TrimSpace(configuredDefault)

	if mb := firstExactMailbox(mailboxes, "INBOX"); mb != "" {
		return mb
	}

	if configured != "" {
		if mb := firstExactMailbox(mailboxes, configured); mb != "" {
			return mb
		}
	}

	for _, mb := range mailboxes {
		if strings.HasPrefix(mb, "---") {
			continue
		}
		if strings.TrimSpace(mb) == "" {
			continue
		}
		return mb
	}

	return ""
}

func firstExactMailbox(mailboxes []string, name string) string {
	for _, mb := range mailboxes {
		if mb == name {
			return mb
		}
	}
	return ""
}
