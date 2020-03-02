# Tracing

Visualize the Requests Flow
{: .subtitle }

The tracing system allows developers to visualize call flows in their infrastructure.

Traefik uses OpenTracing, an open standard designed for distributed tracing.

Traefik supports five tracing backends:

- [Jaeger](./jaeger.md)
- [Zipkin](./zipkin.md)
- [Datadog](./datadog.md)
- [Instana](./instana.md)
- [Haystack](./haystack.md)

## Configuration

By default, Traefik uses Jaeger as tracing backend.

To enable the tracing:

```toml tab="File (TOML)"
[tracing]
```

```yaml tab="File (YAML)"
tracing: {}
```

```bash tab="CLI"
--tracing=true
```

### Common Options

#### `serviceName`

_Required, Default="traefik"_

Service name used in selected backend.

```toml tab="File (TOML)"
[tracing]
  serviceName = "traefik"
```

```yaml tab="File (YAML)"
tracing:
  serviceName: traefik
```

```bash tab="CLI"
--tracing.serviceName=traefik
```

#### `spanNameLimit`

_Required, Default=0_

Span name limit allows for name truncation in case of very long names.
This can prevent certain tracing providers to drop traces that exceed their length limits.

`0` means no truncation will occur.

```toml tab="File (TOML)"
[tracing]
  spanNameLimit = 150
```

```yaml tab="File (YAML)"
tracing:
  spanNameLimit: 150
```

```bash tab="CLI"
--tracing.spanNameLimit=150
```
