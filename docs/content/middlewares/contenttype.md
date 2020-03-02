
# ContentType

Handling ContentType auto-detection
{: .subtitle }

The Content-Type middleware - or rather its unique `autoDetect` option -
specifies whether to let the `Content-Type` header,
if it has not been set by the backend,
be automatically set to a value derived from the contents of the response.

As a proxy, the default behavior should be to leave the header alone,
regardless of what the backend did with it.
However, the historic default was to always auto-detect and set the header if it was nil,
and it is going to be kept that way in order to support users currently relying on it.
This middleware exists to enable the correct behavior until at least the default one can be changed in a future version.

!!! info

    As explained above, for compatibility reasons the default behavior on a router (without this middleware),
    is still to automatically set the `Content-Type` header.
    Therefore, given the default value of the `autoDetect` option (false),
    simply enabling this middleware for a router switches the router's behavior.

## Configuration Examples

```yaml tab="Docker"
# Disable auto-detection
labels:
  - "traefik.http.middlewares.autodetect.contenttype.autodetect=false"
```

```yaml tab="Kubernetes"
# Disable auto-detection
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: autodetect
spec:
  contentType:
    autoDetect: false
```

```yaml tab="Consul Catalog"
# Disable auto-detection
- "traefik.http.middlewares.autodetect.contenttype.autodetect=false"
```

```json tab="Marathon"
"labels": {
  "traefik.http.middlewares.autodetect.contenttype.autodetect": "false"
}
```

```yaml tab="Rancher"
# Disable auto-detection
labels:
  - "traefik.http.middlewares.autodetect.contenttype.autodetect=false"
```

```toml tab="File (TOML)"
# Disable auto-detection
[http.middlewares]
  [http.middlewares.autodetect.contentType]
     autoDetect=false
```

```yaml tab="File (YAML)"
# Disable auto-detection
http:
  middlewares:
    autodetect:
      contentType:
        autoDetect: false
```

## Configuration Options

### `autoDetect`

`autoDetect` specifies whether to let the `Content-Type` header,
if it has not been set by the backend,
be automatically set to a value derived from the contents of the response.
