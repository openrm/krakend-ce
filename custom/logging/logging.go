package logging

import (
	"fmt"
	"time"
	"bytes"
	"regexp"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
)

const (
	Namespace = "github_com/openrm/krakend-logging"
	Prefix = "logging:"
)

type Config struct {
	HeaderBlacklist []string `json:"header_blacklist"`
}

func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}
	cfg := Config{}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&v); err != nil {
		return nil
	}
	if err := json.NewDecoder(buf).Decode(&cfg); err != nil {
		return nil
	}
	return cfg
}

func filterHeaders(headers http.Header, blacklist []*regexp.Regexp) map[string]string {
	filtered := make(map[string]string)
	for k := range headers {
		var value string = headers.Get(k)
		for _, re := range blacklist {
			if re.MatchString(k) {
				value = "[Filtered]"
			}
		}
		filtered[k] = value
	}
	return filtered
}

func GetFormatter(cfg config.ServiceConfig) gin.LogFormatter {
	var (
		blacklist []*regexp.Regexp
		formatter gin.LogFormatter
	)
	if cfg, ok := ConfigGetter(cfg.ExtraConfig).(Config); ok {
		for _, v := range cfg.HeaderBlacklist {
			if exp, err := regexp.Compile(v); err == nil {
				blacklist = append(blacklist, exp)
			}
		}
		formatter = func(param gin.LogFormatterParams) string {
			r := param.Request
			msg := fmt.Sprintf(
				"[GIN] %v | %3d | %13v | %15s | %-7s  %#v\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
			fields := map[string]interface{}{
				"@timestamp": param.TimeStamp.Format(time.RFC3339),
				"level": "info",
				"message": msg,
				"module": "GIN",
				"ip": param.ClientIP,
				"method": param.Method,
				"protocol": r.Proto,
				"url": r.RequestURI,
				"remoteAddress": r.RemoteAddr,
				"hostname": r.Host,
				"referer": r.Referer(),
				"userAgent": r.UserAgent(),
				"contentLength": r.ContentLength,
				"headers": filterHeaders(r.Header, blacklist),
				"responseTime": float64(param.Latency.Nanoseconds()) / 1e6,
				"status": param.StatusCode,
				"responseContentLength": param.BodySize,
			}
			if len(param.ErrorMessage) > 0 {
				fields["err"] = param.ErrorMessage
			}
			bs, _ := json.Marshal(fields)
			return string(bs) + "\n"
		}
	}
	return formatter
}
