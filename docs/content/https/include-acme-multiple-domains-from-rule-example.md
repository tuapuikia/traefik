
```yaml tab="Docker"
## Dynamic configuration
labels:
  - traefik.http.routers.blog.rule=(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)
  - traefik.http.routers.blog.tls=true
  - traefik.http.routers.blog.tls.certresolver=myresolver
```

```yaml tab="Docker (Swarm)"
## Dynamic configuration
deploy:
  labels:
    - traefik.http.routers.blog.rule=(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)
    - traefik.http.routers.blog.tls=true
    - traefik.http.routers.blog.tls.certresolver=myresolver
    - traefik.http.services.blog-svc.loadbalancer.server.port=8080"
```

```yaml tab="Kubernetes"
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: blogtls
spec:
  entryPoints:
    - websecure
  routes:
  - match: (Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)
    kind: Rule
    services:
    - name: blog
      port: 8080
  tls:
    certresolver: myresolver
```

```json tab="Marathon"
labels: {
  "traefik.http.routers.blog.rule": "(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)",
  "traefik.http.routers.blog.tls": "true",
  "traefik.http.routers.blog.tls.certresolver": "myresolver",
  "traefik.http.services.blog-svc.loadbalancer.server.port": "8080"
}
```

```yaml tab="Rancher"
## Dynamic configuration
labels:
  - traefik.http.routers.blog.rule=(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)
  - traefik.http.routers.blog.tls=true
  - traefik.http.routers.blog.tls.certresolver=myresolver
```

```toml tab="File (TOML)"
## Dynamic configuration
[http.routers]
  [http.routers.blog]
    rule = "(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)"
    [http.routers.blog.tls]
      certResolver = "myresolver"
```

```yaml tab="File (YAML)"
## Dynamic configuration
http:
  routers:
    blog:
      rule: "(Host(`company.com`) && Path(`/blog`)) || Host(`blog.company.org`)"
      tls:
        certResolver: myresolver
```
