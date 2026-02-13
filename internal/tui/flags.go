package tui

import "github.com/chhlga/budge/internal/email"

func addFlag(flags []string, flag string) []string {
	for _, f := range flags {
		if f == flag {
			return flags
		}
	}
	return append(flags, flag)
}

func removeFlag(flags []string, flag string) []string {
	out := make([]string, 0, len(flags))
	for _, f := range flags {
		if f != flag {
			out = append(out, f)
		}
	}
	return out
}

func markSeenInSlice(emails []email.Message, uid uint32, seen bool) []email.Message {
	for i := range emails {
		if emails[i].UID != uid {
			continue
		}
		if seen {
			emails[i].Flags = addFlag(emails[i].Flags, "\\Seen")
		} else {
			emails[i].Flags = removeFlag(emails[i].Flags, "\\Seen")
		}
		break
	}

	return emails
}
