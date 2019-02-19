package util

import (
	"strings"
	"github.com/satori/go.uuid"
)

const DASH = "-"

// GetUUID generates a UUID/GUID (version 1 format) using the satori/go.uuid package
func GetUUID() string {
	// Version 1, based on timestamp and MAC address (RFC 4122)
	u1, _ := uuid.NewV1()
	return strings.Replace(u1.String(), string(DASH), "", -1)
}

// GetUUIDv4 generates a UUID/GUID (version 4 format) using the satori/go.uuid package
func GetUUIDv4() string {
	// Version 4, based on random numbers (RFC 4122)
	u4, _ := uuid.NewV4()
	return strings.Replace(u4.String(), string(DASH), "", -1)
}