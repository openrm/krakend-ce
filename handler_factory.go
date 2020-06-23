package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	"github.com/openrm/krakend-jose"
	lua "github.com/devopsfaith/krakend-lua/router/gin"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, lcfg loggingConfig, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	router.RegisterRender("json_error", jsonErrorRender)
	handlerFactory := juju.NewRateLimiterMw(router.EndpointHandler)
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = NewJoseHandlerFactory(handlerFactory, logger, rejecter)
	// handlerFactory = NewRefreshHandlerFactory(handlerFactory, logger)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = NewOpenCensusHandlerFactory(handlerFactory, lcfg)
	handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}
