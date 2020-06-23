package krakend

import (
	"sync"
	"context"
	"net/http"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/gin-gonic/gin"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/devopsfaith/krakend/transport/http/client"
	"github.com/devopsfaith/krakend-opencensus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var once sync.Once

func NewOpenCensusClient(lcfg loggingConfig, clientFactory client.HTTPClientFactory) client.HTTPClientFactory {
	if !lcfg.configured {
		return clientFactory
	}

	prop := lcfg.httpFormat()

	return func(ctx context.Context) *http.Client {
		client := clientFactory(ctx)

		once.Do(func() {
			transport := client.Transport

			if transport == nil {
				transport = http.DefaultTransport
			}

			client.Transport = &ochttp.Transport{
				Base: transport,
				Propagation: prop,
			}
		})

		return client
	}
}

func NewOpenCensusHandlerFactory(hf router.HandlerFactory, lcfg loggingConfig) router.HandlerFactory {
	if !lcfg.configured {
		return hf
	}

	skip, prop := lcfg.skipPaths, lcfg.httpFormat()
	filterPath := func(r *http.Request) trace.StartOptions{
		if u := r.URL; u != nil {
			if _, ok := skip[u.Path]; ok {
				return trace.StartOptions{
					Sampler: trace.NeverSample(),
				}
			}
		}
		return trace.StartOptions{}
	}

	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handler := hf(cfg, p)
		return func(c *gin.Context) {
			traceHandler := ochttp.Handler{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Set(opencensus.ContextKey, trace.FromContext(r.Context()))
					c.Request = r
					handler(c)
				}),
				Propagation: prop,
				GetStartOptions: filterPath,
				FormatSpanName: func(*http.Request) string {
					return cfg.Endpoint
				},
			}
			traceHandler.ServeHTTP(c.Writer, c.Request)
		}
	}
}
