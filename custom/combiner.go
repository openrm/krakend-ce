package custom

import (
	"net/http"
	"github.com/luraproject/lura/proxy"
)

func init() {
	proxy.RegisterResponseCombiner("status", combineStatusAndData)
}

// Custom combiner to store the backend response's status code < 500 in the metadata,
// but only when there is exactly one backend or the codes are uniform.
// Can be used together with the custom renderer that forwards this status code
func combineStatusAndData(total int, parts []*proxy.Response) *proxy.Response {
	isComplete := len(parts) == total

	var (
		retResponse *proxy.Response
		retStatus int
	)

	for _, part := range parts {
		if part == nil || part.Data == nil {
			isComplete = false
			continue
		}

		if status := part.Metadata.StatusCode; status != 0 { // failed
			if retStatus == 0 {
				retStatus = status
			} else if retStatus != status { // at least two status codes are different
				retStatus = http.StatusInternalServerError
			}
		}

		isComplete = isComplete && part.IsComplete
		if retResponse == nil {
			retResponse = part
			continue
		}

		for k, v := range part.Data {
			retResponse.Data[k] = v
		}
	}

	if retStatus > http.StatusInternalServerError {
		retStatus = http.StatusInternalServerError
	}

	if retResponse == nil {
		// do not allow nil data in the response:
		return &proxy.Response{
			Data: make(map[string]interface{}, 0),
			IsComplete: isComplete,
			Metadata: proxy.Metadata{StatusCode: retStatus},
		}
	}

	retResponse.IsComplete = isComplete
	retResponse.Metadata.StatusCode = retStatus
	return retResponse
}
