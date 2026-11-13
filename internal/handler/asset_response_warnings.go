package handler

import "context"

const (
	assetHeaderConflictWarning = "asset metadata header conflicts with storage header"
	assetUnknownLengthWarning  = "asset metadata content length ignored because storage length is unknown"
)

type assetHeaderResponseContext struct {
	warn     func(context.Context, string, ...any)
	ctx      context.Context
	filename string
}

func (r assetHeaderResponseContext) warnConflict(name, source string) {
	r.warnHeader(assetHeaderConflictWarning, name, source)
}

func (r assetHeaderResponseContext) warnUnknownContentLength(name string) {
	r.warnHeader(assetUnknownLengthWarning, name, "storage")
}

func (r assetHeaderResponseContext) warnHeader(message, name, source string) {
	if r.warn == nil {
		return
	}
	r.warn(
		r.ctx,
		message,
		"filename", r.filename,
		"header", name,
		"source", source,
	)
}
