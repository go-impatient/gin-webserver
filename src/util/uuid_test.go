package util

import (
	"strings"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	if len(uuid) == 0 || strings.Contains(uuid, "-") {
		t.Fatalf("TestGenerateUuid failed")
	}
}
