# Instana Exporter

| Status                   |                  |
| ------------------------ |------------------|
| Stability                | [beta]           |
| Supported pipeline types | traces           |
| Distributions            | [contrib], [AWS] |

The Instana Exporter takes the role of the Instana Agent at exporting observability signals to the Instana platform.

The Instana Agent is already able to receive otlp spans, converts them into Instana specific spans and send them to the Instana Backend.
However, customers who wish to adopt the OpenTelemetry Collector can replace the Instana Agent by their custom OpenTelemetry Collector, by including the Instana Exporter into the build.


## Configuration

In order to use the Instana Exporter within your Collector, two steps must be followed:

1. Add the ``instana`` entry under the ``exporters`` section in your configuration file.
1. Add the ``instana`` entry as part of pipelines under the ``service/pipelines/traces`` entry.

### Parameters

The Instana Exporter requires three parameters in order to startup properly and be able to send spans to the Instana Backend:

 * ``agent_endpoint``: The Instana backend endpoint that the Exporter connects to. It depends on your region and it starts with ``https://serverless-``. It corresponds to the Instana environment variable ``INSTANA_ENDPOINT_URL``.
 * ``agent_key``: Your Instana Agent key. The same agent key can be used for host agents and serverless monitoring. It corresponds to the Instana environment variable ``INSTANA_AGENT_KEY``.
 * ``zone``: The zone to place this monitored component into. It corresponds to the Instana environment variable ``INSTANA_ZONE``.

> These parameters match the Instana Serverless Monitoring environment variables and can be found [here](https://www.ibm.com/docs/en/instana-observability/current?topic=references-environment-variables#serverless-monitoring).

### Sample Configuration

The code snippet below shows how your configuration file should look like:

```yaml
[...]

exporters:
  instana:
    agent_endpoint: ${INSTANA_ENDPOINT_URL}
    agent_key: ${INSTANA_AGENT_KEY}
    zone: ${INSTANA_ZONE}

[...]

service:
  pipelines:
    traces:
      exporters: [logging, instana]

[...]
```

### Full Example

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
exporters:
  logging:
    loglevel: debug
  instana:
    loglevel: debug
    agent_endpoint: ${INSTANA_ENDPOINT_URL}
    agent_key: ${INSTANA_AGENT_KEY}
    zone: ${INSTANA_ZONE}

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, instana]
    metrics:
      receivers: [otlp]
      processors: [ batch]
      exporters: [logging]
```


[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[AWS]:https://aws-otel.github.io/docs/partners/dynatrace
