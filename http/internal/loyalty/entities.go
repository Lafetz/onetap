package core

import "github.com/google/uuid"

type Notification struct {
	CustomerID uuid.UUID
	Desc       string
}
