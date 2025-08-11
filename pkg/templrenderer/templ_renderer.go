package templrenderer

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/jfyne/live"
)

type LiveTempl func(rc *live.RenderContext) templ.Component

// type RenderHandler func(ctx context.Context, rc *live.RenderContext) (io.Reader, error)

// WithTemplateRenderer set the handler to use an `html/template` renderer.
func WithTemplRenderer(lt LiveTempl) live.HandlerConfig {
	return func(h *live.Handler) error {
		h.RenderHandler = func(ctx context.Context, rc *live.RenderContext) (io.Reader, error) {
			var buf bytes.Buffer
			if err := lt(rc).Render(ctx, &buf); err != nil {
				return nil, err
			}
			return &buf, nil
		}
		return nil
	}
}
