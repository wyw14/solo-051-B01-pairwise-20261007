package model

import "github.com/google/uuid"

func NewID() ID { return ID(uuid.NewString()) }

func ParseID(value string) (ID, bool) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return "", false
	}
	return ID(parsed.String()), true
}
