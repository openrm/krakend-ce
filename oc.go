package krakend

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/config"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	orprop "github.com/openrm/module-tracing-golang/propagation"
)

func NewOpenCensusMiddleware(cfg config.ExtraConfig) gin.HandlerFunc {
	var prop propagation.HTTPFormat
	var skip map[string]struct{}
	if v, ok := cfg[LoggerNamespace]; ok {
		if cfg, ok := v.(map[string]interface{}); ok {
			if v, ok := cfg["skip_paths"]; ok {
				if ps, ok := v.([]interface{}); ok && len(ps) > 0 {
					skip = make(map[string]struct{}, len(ps))
					for _, v := range ps {
						if path, ok := v.(string); ok {
							skip[path] = struct{}{}
						}
					}
				}
			}
			if v, ok := cfg["trace_header"]; ok {
				if h, ok := v.(string); ok {
					prop = &orprop.HTTPFormat{Header: h}
				}
			}
		}
	}
	handler := ochttp.Handler{
		Propagation: prop,
		GetStartOptions: func(r *http.Request) trace.StartOptions{
			if u := r.URL; u != nil {
				if _, ok := skip[u.Path]; ok {
					return trace.StartOptions{
						Sampler: trace.NeverSample(),
					}
				}
			}
			return trace.StartOptions{}
		},
	}
	return func(c *gin.Context) {
		handler.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
