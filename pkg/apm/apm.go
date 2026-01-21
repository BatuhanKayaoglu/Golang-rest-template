package apm

import (
	"context"
	"log"
	"os"

	"golang-rest-api-template/pkg/config"

	"go.elastic.co/apm/v2"
)

// Init initializes the APM tracer with configuration from environment variables.
// The Elastic APM Go agent automatically reads configuration from environment variables,
// but we set them explicitly here to ensure they match our config.
func Init(cfg *config.APMConfig) {
	if !cfg.Active {
		log.Println("APM is disabled")
		return
	}

	// Set environment variables that the APM agent reads
	os.Setenv("ELASTIC_APM_SERVICE_NAME", cfg.ServiceName)
	os.Setenv("ELASTIC_APM_SERVER_URL", cfg.ServerURL)
	os.Setenv("ELASTIC_APM_ENVIRONMENT", cfg.Environment)
	if cfg.SecretToken != "" {
		os.Setenv("ELASTIC_APM_SECRET_TOKEN", cfg.SecretToken)
	}

	log.Printf("APM initialized: service=%s, server=%s, environment=%s",
		cfg.ServiceName, cfg.ServerURL, cfg.Environment)
}

// Close closes the APM tracer and flushes any pending data.
func Close() {
	tracer := apm.DefaultTracer()
	tracer.Flush(nil)
	tracer.Close()
	log.Println("APM tracer closed")
}

// CaptureError captures an error and sends it to APM.
func CaptureError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	apm.CaptureError(ctx, err).Send()
}

// StartSpan starts a new span with the given name and type.
// Returns the span and a context with the span.
func StartSpan(ctx context.Context, name, spanType string) (*apm.Span, context.Context) {
	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		return nil, ctx
	}
	span := tx.StartSpan(name, spanType, nil)
	return span, apm.ContextWithSpan(ctx, span)
}

// EndSpan ends the given span.
func EndSpan(span *apm.Span) {
	if span != nil {
		span.End()
	}
}

// SetTransactionResult sets the result of the current transaction.
func SetTransactionResult(ctx context.Context, result string) {
	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		tx.Result = result
	}
}

// AddLabel adds a label to the current transaction or span.
func AddLabel(ctx context.Context, key string, value interface{}) {
	if tx := apm.TransactionFromContext(ctx); tx != nil {
		tx.Context.SetLabel(key, value)
	}
}
