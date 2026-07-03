<!-- togo-brand -->
<p align="center">
  <img src=".github/assets/togo-mark.svg" width="96" alt="togo" />
</p>
<h1 align="center">providers</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a> — the full-stack Go + React framework</sub></p>

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

MIT © fadymondy
