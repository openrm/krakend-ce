package custom

import (
	"github.com/openrm/module-tracing-golang/propagation"
)

const traceHeader = "X-Openrm-Trace"

var Propagation = &propagation.HTTPFormat{Header: traceHeader}
