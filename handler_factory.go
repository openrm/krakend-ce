package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	jose "github.com/devopsfaith/krakend-jose"
	ginjose "github.com/devopsfaith/krakend-jose/gin"
	lua "github.com/devopsfaith/krakend-lua/router/gin"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus/router/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/luraproject/lura/logging"
	router "github.com/luraproject/lura/router/gin"
)

import (
	sentry "github.com/openrm/krakend-sentry/gin"
	"github.com/devopsfaith/krakend-ce/custom"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	handlerFactory := juju.HandlerFactory
	handlerFactory = sentry.HandlerFactory(logger, handlerFactory)
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter, custom.StatusRejecterFactory)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.NewWithPropagation(handlerFactory, custom.Propagation)
	handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}

type handlerFactory struct{}

func (h handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory) router.HandlerFactory {
	return NewHandlerFactory(l, m, r)
}
