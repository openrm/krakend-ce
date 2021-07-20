package custom

import (
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	cel "github.com/devopsfaith/krakend-cel"
	jose "github.com/devopsfaith/krakend-jose"
)

var StatusRejecterFactory = jose.StatusRejecterFactoryFunc(func(l logging.Logger, cfg *config.EndpointConfig) jose.StatusRejecter {
	if r := cel.NewStatusRejecter(l, cfg); r != nil {
		return r
	}
	return jose.FixedStatusRejecter{false, 0}
})
