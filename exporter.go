package instanaexporter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/ptrace"

	instanaConfig "github.com/ibm-observability/instanaexporter/config"
	"github.com/ibm-observability/instanaexporter/internal/converter"
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	"github.com/ibm-observability/instanaexporter/internal/otlptext"
)

type instanaExporter struct {
	config          *instanaConfig.Config
	client          *http.Client
	logger          *zap.Logger
	tracesMarshaler ptrace.Marshaler
	settings        component.TelemetrySettings
	userAgent       string
}

func (e *instanaExporter) start(_ context.Context, host component.Host) error {
	client, err := e.config.HTTPClientSettings.ToClient(host, e.settings)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

func (e *instanaExporter) pushConvertedTraces(ctx context.Context, td ptrace.Traces) error {
	e.logger.Info("TracesExporter", zap.Int("#spans", td.SpanCount()))
	if !e.logger.Core().Enabled(zapcore.DebugLevel) {
		return nil
	}

	buf, err := e.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	e.logger.Debug(string(buf))

	converter := converter.NewConvertAllConverter(e.logger)
	spans := make([]model.Span, 0)

	hostId := ""
	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resSpan := resourceSpans.At(i)

		resource := resSpan.Resource()

		hostIdAttr, ex := resource.Attributes().Get(instanaConfig.AttributeInstanaHostID)
		if ex {
			hostId = hostIdAttr.StringVal()
		}

		ilSpans := resSpan.ScopeSpans()
		for j := 0; j < ilSpans.Len(); j++ {
			converterBundle := converter.ConvertSpans(resource.Attributes(), ilSpans.At(j).Spans())

			spans = append(spans, converterBundle.Spans...)
		}
	}

	bundle := model.Bundle{Spans: spans}

	if len(bundle.Spans) <= 0 {
		// skip exporting, nothing to do
		return nil
	}

	req, err := bundle.Marshal()

	e.logger.Debug(string(req))

	if err != nil {
		return consumererror.NewPermanent(err)
	}

	headers := map[string]string{
		instanaConfig.HeaderKey:  e.config.AgentKey,
		instanaConfig.HeaderHost: hostId,
		instanaConfig.HeaderTime: "0",
	}

	return e.export(ctx, e.config.Endpoint, headers, req)
}

func newInstanaExporter(logger *zap.Logger, cfg config.Exporter, set component.ExporterCreateSettings) (*instanaExporter, error) {
	iCfg := cfg.(*instanaConfig.Config)

	if iCfg.Endpoint != "" {
		_, err := url.Parse(iCfg.Endpoint)
		if err != nil {
			return nil, errors.New("endpoint must be a valid URL")
		}
	}

	userAgent := fmt.Sprintf("%s/%s (%s/%s)", set.BuildInfo.Description, set.BuildInfo.Version, runtime.GOOS, runtime.GOARCH)

	return &instanaExporter{
		config:          iCfg,
		logger:          logger,
		tracesMarshaler: otlptext.NewTextTracesMarshaler(),
		userAgent:       userAgent,
	}, nil
}

func (e *instanaExporter) export(ctx context.Context, url string, header map[string]string, request []byte) error {
	url = strings.TrimSuffix(url, "/") + "/bundle"

	e.logger.Debug("Preparing to make HTTP request", zap.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(request))
	if err != nil {
		return consumererror.NewPermanent(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", e.userAgent)

	for name, value := range header {
		req.Header.Set(name, value)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make an HTTP request: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Request is successful.
		return nil
	}

	return nil
}
