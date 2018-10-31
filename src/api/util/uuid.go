package util

import (
	"github.com/satori/go.uuid"
	"strings"
)

const DASH = "-"

func GenerateUuid() string {
	return strings.Replace(uuid.NewV1().String(), string(DASH), "", -1)
}

