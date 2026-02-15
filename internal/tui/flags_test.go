package tui

import (
	"reflect"
	"testing"

	"github.com/chhlga/budge/internal/email"
)

func TestAddFlag(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		flag     string
		expected []string
	}{
		{
			name:     "add to empty",
			flags:    []string{},
			flag:     "\\Seen",
			expected: []string{"\\Seen"},
		},
		{
			name:     "add to existing",
			flags:    []string{"\\Flagged"},
			flag:     "\\Seen",
			expected: []string{"\\Flagged", "\\Seen"},
		},
		{
			name:     "add duplicate",
			flags:    []string{"\\Seen", "\\Flagged"},
			flag:     "\\Seen",
			expected: []string{"\\Seen", "\\Flagged"},
		},
		{
			name:     "add to multiple",
			flags:    []string{"\\Seen", "\\Flagged", "\\Draft"},
			flag:     "\\Answered",
			expected: []string{"\\Seen", "\\Flagged", "\\Draft", "\\Answered"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addFlag(tt.flags, tt.flag)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("addFlag(%v, %q) = %v, want %v", tt.flags, tt.flag, result, tt.expected)
			}
		})
	}
}

func TestRemoveFlag(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		flag     string
		expected []string
	}{
		{
			name:     "remove from empty",
			flags:    []string{},
			flag:     "\\Seen",
			expected: []string{},
		},
		{
			name:     "remove existing",
			flags:    []string{"\\Seen", "\\Flagged"},
			flag:     "\\Seen",
			expected: []string{"\\Flagged"},
		},
		{
			name:     "remove non-existing",
			flags:    []string{"\\Flagged"},
			flag:     "\\Seen",
			expected: []string{"\\Flagged"},
		},
		{
			name:     "remove from single",
			flags:    []string{"\\Seen"},
			flag:     "\\Seen",
			expected: []string{},
		},
		{
			name:     "remove middle flag",
			flags:    []string{"\\Seen", "\\Flagged", "\\Draft"},
			flag:     "\\Flagged",
			expected: []string{"\\Seen", "\\Draft"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeFlag(tt.flags, tt.flag)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("removeFlag(%v, %q) = %v, want %v", tt.flags, tt.flag, result, tt.expected)
			}
		})
	}
}

func TestMarkSeenInSlice(t *testing.T) {
	tests := []struct {
		name     string
		emails   []email.Message
		uid      uint32
		seen     bool
		expected []email.Message
	}{
		{
			name:     "mark seen in empty slice",
			emails:   []email.Message{},
			uid:      1,
			seen:     true,
			expected: []email.Message{},
		},
		{
			name: "mark message as seen",
			emails: []email.Message{
				{UID: 1, Flags: []string{}},
				{UID: 2, Flags: []string{"\\Flagged"}},
			},
			uid:  1,
			seen: true,
			expected: []email.Message{
				{UID: 1, Flags: []string{"\\Seen"}},
				{UID: 2, Flags: []string{"\\Flagged"}},
			},
		},
		{
			name: "mark message as unseen",
			emails: []email.Message{
				{UID: 1, Flags: []string{"\\Seen"}},
				{UID: 2, Flags: []string{"\\Flagged"}},
			},
			uid:  1,
			seen: false,
			expected: []email.Message{
				{UID: 1, Flags: []string{}},
				{UID: 2, Flags: []string{"\\Flagged"}},
			},
		},
		{
			name: "mark non-existent uid",
			emails: []email.Message{
				{UID: 1, Flags: []string{}},
				{UID: 2, Flags: []string{}},
			},
			uid:  999,
			seen: true,
			expected: []email.Message{
				{UID: 1, Flags: []string{}},
				{UID: 2, Flags: []string{}},
			},
		},
		{
			name: "mark already seen message as seen",
			emails: []email.Message{
				{UID: 1, Flags: []string{"\\Seen"}},
			},
			uid:  1,
			seen: true,
			expected: []email.Message{
				{UID: 1, Flags: []string{"\\Seen"}},
			},
		},
		{
			name: "mark unseen message as unseen",
			emails: []email.Message{
				{UID: 1, Flags: []string{"\\Flagged"}},
			},
			uid:  1,
			seen: false,
			expected: []email.Message{
				{UID: 1, Flags: []string{"\\Flagged"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := markSeenInSlice(tt.emails, tt.uid, tt.seen)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("markSeenInSlice() mismatch:\ngot:  %+v\nwant: %+v", result, tt.expected)
			}
		})
	}
}
