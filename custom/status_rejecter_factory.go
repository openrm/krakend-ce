package custom

import (
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	cel "github.com/krakendio/krakend-cel/v2"
	jose "github.com/krakendio/krakend-jose/v2"
)

var StatusRejecterFactory = jose.StatusRejecterFactoryFunc(func(l logging.Logger, cfg *config.EndpointConfig) jose.StatusRejecter {
	if r := cel.NewStatusRejecter(l, cfg); r != nil {
		return r
	}
	return jose.FixedStatusRejecter{false, 0}
})
