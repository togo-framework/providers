<!-- togo-header -->
<div align="center">
  <picture><source media="(prefers-color-scheme: dark)" srcset=".github/assets/togo-mark-dark.svg" /><img src=".github/assets/togo-mark.svg" alt="ToGO" height="64" /></picture>
  <h1>togo-framework/providers</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1F8A99" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/providers"><img src="https://pkg.go.dev/badge/github.com/togo-framework/providers.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/providers
```

<!-- /togo-header -->

Runtime, config-driven **selection** of which backend is active for a capability —
swap it over **CLI / `.env` / settings** without recompiling.

togo already has the registry (`RegisterProviderFunc`), the capability contracts
(`togo.Cache` / `Queue` / `Storage` / `Broker`), the container (`k.Set/Get`), and
the `<capability>-<backend>` plugin convention (`storage-s3`, `realtime-nats`, …).
The one missing piece is picking **which** backend wins at runtime. This package
adds exactly that — and nothing else. It stores no config itself; the
[`settings`](https://github.com/togo-framework/settings) plugin is the store.

## Selection order (highest wins)

```
1. env        TOGO_<CAP>_PROVIDER=iceberg      # or a CLI flag that sets it
2. settings   providers.<cap>.active           # togo.yaml / settings plugin
3. default    the backend registered with def=true
```

## Use it — in a backend plugin

```go
import "github.com/togo-framework/providers"

func init() {
  togo.RegisterProviderFunc("db-iceberg", togo.PriorityService, func(k *togo.Kernel) error {
    // installs itself as the "data" service ONLY if selected (def=false → opt-in)
    providers.Use(k, providers.CapData, "iceberg", newIceberg(k), false)
    return nil
  })
}
```

Read provider config the dynamic way (env → settings → default; secrets env-only):

```go
url   := providers.Value(k, providers.CapExecute, "coder", "url",   "http://localhost:7080", false)
token := providers.Value(k, providers.CapExecute, "coder", "token", "",                      true) // secret
```

## Capabilities

| Constant | Contract | Backends |
|---|---|---|
| `CapCache` `CapQueue` `CapStorage` `CapRealtime` | **existing** (`togo` core) | e.g. `storage-s3`, `queue-nats`, `realtime-nats` |
| `CapImplement` `CapExecute` | autopilot | `claude`/`omnigent`, `local`/`coder` |
| `CapCompute` | worker | `worker-local`/`beam`/`spark`/`flink`/`databricks` |
| `CapData` `CapCatalog` `CapShare` | db / catalog / sharing | `pg`/`iceberg`/`bigquery`/`databricks`, … |

## Install

```bash
togo install togo-framework/providers
```
<!-- togo-sponsors -->
---

<div align="center">
  <h3>💎 Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><img src=".github/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
    &nbsp;&nbsp;&nbsp;&nbsp;
    <a href="https://one-studio.co"><img src=".github/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
