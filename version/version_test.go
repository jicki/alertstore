package version_test

import (
	"testing"

	"alertstort/version"
)

func TestVersionString(t *testing.T) {
	version.Commit = "mycommit"
	version.Date = "today"
	version.Version = "0.0.1"

	expected := "alertsnitch Version: 0.0.1 Commit: mycommit Date: today"
	if version.GetVersion() != expected {
		t.Fatalf("invalid version %s expected %s", version.GetVersion(), expected)
	}
}
