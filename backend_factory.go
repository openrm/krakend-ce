package krakend

import (
	"net/http"
	"context"

	// amqp "github.com/devopsfaith/krakend-amqp"
	cel "github.com/devopsfaith/krakend-cel"
	cb "github.com/devopsfaith/krakend-circuitbreaker/gobreaker/proxy"
	httpcache "github.com/devopsfaith/krakend-httpcache"
	lambda "github.com/devopsfaith/krakend-lambda"
	lua "github.com/devopsfaith/krakend-lua/proxy"
	"github.com/devopsfaith/krakend-martian"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend-oauth2-clientcredentials"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	pubsub "github.com/devopsfaith/krakend-pubsub"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/proxy"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/encoding"
	"github.com/devopsfaith/krakend/transport/http/client"
	"go.opencensus.io/trace"
	httprequestexecutor "github.com/devopsfaith/krakend/transport/http/client/plugin"
)

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares:
// - oauth2 client credentials
// - http cache
// - martian
// - pubsub
// - amqp
// - cel
// - lua
// - rate-limit
// - circuit breaker
// - metrics collector
// - opencensus collector
// func NewBackendFactory(logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
// 	return NewBackendFactoryWithContext(context.Background(), logger, metricCollector)
// }

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares and injecting the received context
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, lcfg loggingConfig, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		var clientFactory client.HTTPClientFactory
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		} else {
			clientFactory = httpcache.NewHTTPClient(cfg)
		}
		clientFactory = NewOpenCensusClient(lcfg, clientFactory)
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			re := opencensus.HTTPRequestExecutor(clientFactory)
			return re(trace.NewContext(ctx, ctx.Value(opencensus.ContextKey).(*trace.Span)), req)
		}
	}
	requestExecutorFactory = httprequestexecutor.HTTPRequestExecutor(logger, requestExecutorFactory)
	//  the line below registers martian.staticModifierFromJSON
	var _ = martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)

	backendFactory := func(cfg *config.Backend) proxy.Proxy {
		re := requestExecutorFactory(cfg)
		if result, ok := martian.ConfigGetter(cfg.ExtraConfig).(martian.Result); ok {
			if result.Err == nil {
				re = martian.HTTPRequestExecutor(result.Result, re)
			} else if result.Err != martian.ErrEmptyValue {
				logger.Error(result, cfg.ExtraConfig)
			}
		}

		rp := proxy.NoOpHTTPResponseParser
		if cfg.Encoding != encoding.NOOP {
			ef := proxy.NewEntityFormatter(cfg)
			rp = proxy.DefaultHTTPResponseParserFactory(proxy.HTTPResponseParserConfig{cfg.Decoder, ef})
		}

		return proxy.NewHTTPProxyDetailed(cfg, re, GetHTTPStatusHandler(cfg), rp)
	}

	bf := pubsub.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = bf.New
	// backendFactory = amqp.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = lambda.BackendFactory(backendFactory)
	backendFactory = cel.BackendFactory(logger, backendFactory)
	backendFactory = lua.BackendFactory(logger, backendFactory)
	backendFactory = juju.BackendFactory(backendFactory)
	backendFactory = cb.BackendFactory(backendFactory, logger)
	backendFactory = metricCollector.BackendFactory("backend", backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)
	return backendFactory
}
