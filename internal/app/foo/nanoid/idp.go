package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid"

	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
)

var _ service.NanoIdentityProvider = (*uuidIdentityProvider)(nil)

type uuidIdentityProvider struct{}

// New instantiates a UUID identity provider.
func New() service.NanoIdentityProvider {
	return &uuidIdentityProvider{}
}

func (idp *uuidIdentityProvider) ID() (string, error) {
	return gonanoid.Nanoid(21)
}
