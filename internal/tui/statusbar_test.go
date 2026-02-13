package tui

import "testing"

func TestStatusBar_clearsLoadingOnErrorMsg(t *testing.T) {
	sb := NewStatusBar()

	sb, _ = sb.Update(LoadingMsg{Text: "Searching..."})
	if !sb.loading {
		t.Fatalf("expected loading=true")
	}

	sb, _ = sb.Update(ErrorMsg{Err: errTest{}})
	if sb.loading {
		t.Fatalf("expected loading=false")
	}
}

type errTest struct{}

func (errTest) Error() string { return "test" }
