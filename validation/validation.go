package validation

import (
	"link_watcher/serviceErrors"
	"net/url"
)

type TargetRequest struct {
	Url         string `json:"url"`
	IntervalSec int64  `json:"intervalSec"`
}

func ValidateTargetRequest(req *TargetRequest) error {
	if req.Url == "" {
		return serviceErrors.ErrEmptyURL
	}

	if len(req.Url) > 2048 {
		return serviceErrors.ErrURLTooLarge
	}

	parsed, errParse := url.Parse(req.Url)
	if errParse != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return serviceErrors.ErrInvalidURL
	}

	if req.IntervalSec == 0 {
		req.IntervalSec = 60
	}
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
		return serviceErrors.ErrInvalidInterval
	}
	return nil
}
