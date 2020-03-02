# Data Collection

Understanding How Traefik is Being Used
{: .subtitle }

## Configuration Example

Understanding how you use Traefik is very important to us: it helps us improve the solution in many different ways.  
For this very reason, the sendAnonymousUsage option is mandatory: we want you to take time to consider whether or not you wish to share anonymous data with us so we can benefit from your experience and use cases.

!!! example "Enabling Data Collection"
    
    ```toml tab="File (TOML)"
    [global]
      # Send anonymous usage data
      sendAnonymousUsage = true
    ```
    
    ```yaml tab="File (YAML)"
    global:
      # Send anonymous usage data
      sendAnonymousUsage: true
    ```
    
    ```bash tab="CLI"
    # Send anonymous usage data
    --global.sendAnonymousUsage
    ```

## Collected Data

This feature comes from the public proposal [here](https://github.com/containous/traefik/issues/2369).

In order to help us learn more about how Traefik is being used and improve it, we collect anonymous usage statistics from running instances.
Those data help us prioritize our developments and focus on what's important for our users (for example, which provider is popular, and which is not).

### What's collected / when ?

Once a day (the first call begins 10 minutes after the start of Traefik), we collect:

- the Traefik version number
- a hash of the configuration
- an **anonymized version** of the static configuration (token, user name, password, URL, IP, domain, email, etc, are removed).

!!! info
    
    - We do not collect the dynamic configuration information (routers & services).
    - We do not collect this data to run advertising programs.
    - We do not sell this data to third-parties.

### Example of Collected Data

```toml tab="Original configuration"
[entryPoints]
  [entryPoints.web]
    address = ":80"

[api]

[providers.docker]
  endpoint = "tcp://10.10.10.10:2375"
  exposedByDefault = true
  swarmMode = true

  [providers.docker.TLS]
    ca = "dockerCA"
    cert = "dockerCert"
    key = "dockerKey"
    insecureSkipVerify = true
```

```toml tab="Resulting Obfuscated Configuration"
[entryPoints]
  [entryPoints.web]
    address = ":80"

[api]

[providers.docker]
  endpoint = "xxxx"
  exposedByDefault = true
  swarmMode = true

  [providers.docker.TLS]
    ca = "xxxx"
    cert = "xxxx"
    key = "xxxx"
    insecureSkipVerify = false
```

## The Code for Data Collection

If you want to dig into more details, here is the source code of the collecting system: [collector.go](https://github.com/containous/traefik/blob/master/pkg/collector/collector.go)

By default we anonymize all configuration fields, except fields tagged with `export=true`.
