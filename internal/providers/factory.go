package providers

import (
	"fmt"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/httpclient"
)

type Factory struct {
	client *httpclient.Client
}

func NewFactory(client *httpclient.Client) *Factory {
	return &Factory{
		client: client,
	}
}

func (f *Factory) Get(provider string) (Provider, error) {
	switch provider {

	case Greenhouse:
		return NewGreenhouseProvider(f.client), nil

	case Lever:
		return NewLeverProvider(f.client), nil

	case Ashby:
		return NewAshbyProvider(f.client), nil

	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}
