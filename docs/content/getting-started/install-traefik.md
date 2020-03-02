# Install Traefik

You can install Traefik with the following flavors:

* [Use the official Docker image](./#use-the-official-docker-image)
* [(Experimental) Use the Helm Chart](./#use-the-helm-chart)
* [Use the binary distribution](./#use-the-binary-distribution)
* [Compile your binary from the sources](./#compile-your-binary-from-the-sources)

## Use the Official Docker Image

Choose one of the [official Docker images](https://hub.docker.com/_/traefik) and run it with the [sample configuration file](https://raw.githubusercontent.com/containous/traefik/v2.1/traefik.sample.toml):

```bash
docker run -d -p 8080:8080 -p 80:80 \
    -v $PWD/traefik.toml:/etc/traefik/traefik.toml traefik:v2.1
```

For more details, go to the [Docker provider documentation](../providers/docker.md)

!!! tip

    * Prefer a fixed version than the latest that could be an unexpected version.
    ex: `traefik:v2.1.4`
    * Docker images are based from the [Alpine Linux Official image](https://hub.docker.com/_/alpine).
    * Any orchestrator using docker images can fetch the official Traefik docker image.

## Use the Helm Chart

!!! warning "Experimental Helm Chart"
    
    Please note that the Helm Chart for Traefik v2 is still experimental.
    
    The Traefik Stable Chart from 
    [Helm's default charts repository](https://github.com/helm/charts/tree/master/stable/traefik) is still using [Traefik v1.7](https://docs.traefik.io/v1.7).

Traefik can be installed in Kubernetes using the v2.0 Helm chart from <https://github.com/containous/traefik-helm-chart>.

Ensure that the following requirements are met:

* Kubernetes 1.14+
* Helm version 2.x is [installed](https://v2.helm.sh/docs/using_helm/) and initialized with Tiller

Retrieve the latest chart version from the repository:

```bash
# Retrieve Chart from the repository
git clone https://github.com/containous/traefik-helm-chart
```

And install it with the `helm` command line:

```bash
helm install ./traefik-helm-chart
```

!!! tip "Helm Features"
    
    All [Helm features](https://v2.helm.sh/docs/using_helm/#using-helm) are supported.
    For instance, installing the chart in a dedicated namespace:

    ```bash tab="Install in a Dedicated Namespace"
    # Install in the namespace "traefik-v2"
    helm install --namespace=traefik-v2 \
        ./traefik-helm-chart
    ```

??? example "Installing with Custom Values"
    
    You can customize the installation by specifying custom values,
    as with [any helm chart](https://v2.helm.sh/docs/using_helm/#customizing-the-chart-before-installing).
    {: #helm-custom-values }
    
    The values are not (yet) documented, but are self-explanatory:
    you can look at the [default `values.yaml`](https://github.com/containous/traefik-helm-chart/blob/master/traefik/values.yaml) file to explore possibilities.
    
    Example of installation with logging set to `DEBUG`:
    
    ```bash tab="Using Helm CLI"
    helm install --namespace=traefik-v2 \
        --set="logs.loglevel=DEBUG" \
        ./traefik-helm-chart
    ```
    
    ```yml tab="With a custom values file"
    # File custom-values.yml
    ## Install with "helm install --values=./custom-values.yml ./traefik-helm-chart
    logs:
        loglevel: DEBUG
    ```

## Use the Binary Distribution

Grab the latest binary from the [releases](https://github.com/containous/traefik/releases) page.

??? info "Check the integrity of the downloaded file"

    ```bash tab="Linux"
    # Compare this value to the one found in traefik-${traefik_version}_checksums.txt
    sha256sum ./traefik_${traefik_version}_linux_${arch}.tar.gz
    ```

    ```bash tab="macOS"
    # Compare this value to the one found in traefik-${traefik_version}_checksums.txt
    shasum -a256 ./traefik_${traefik_version}_darwin_amd64.tar.gz
    ```

    ```powershell tab="Windows PowerShell"
    # Compare this value to the one found in traefik-${traefik_version}_checksums.txt
    Get-FileHash ./traefik_${traefik_version}_windows_${arch}.zip -Algorithm SHA256
    ```

??? info "Extract the downloaded archive"

    ```bash tab="Linux"
    tar -zxvf traefik_${traefik_version}_linux_${arch}.tar.gz
    ```

    ```bash tab="macOS"
    tar -zxvf ./traefik_${traefik_version}_darwin_amd64.tar.gz
    ```

    ```powershell tab="Windows PowerShell"
    Expand-Archive traefik_${traefik_version}_windows_${arch}.zip
    ```

And run it:

```bash
./traefik --help
```

## Compile your Binary from the Sources

All the details are available in the [Contributing Guide](../contributing/building-testing.md)
