package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cage1016/gae-lab-001/internal/pkg/responses"
)

type Foo struct {
	ID        string    `json:"id" db:"id"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func (p *Foo) MarshalJSON() ([]byte, error) {
	type Alias Foo
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	})
}

type FooItemPage struct {
	responses.Paging
	Items []Foo `json:"items"`
}

// FooRepository specifies an account persistence API.
type FooRepository interface {
	Insert(context.Context, Foo) (string, error)

	RetrieveAll(context.Context, uint64, uint64) (FooItemPage, error)
}
