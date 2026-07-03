# providers

See the [main README](../README.md). Full framework docs: https://github.com/togo-framework/togo/tree/main/docs

Runtime, config-driven **selection** of which backend is active for a capability —
swappable over CLI / `.env` / settings, built on the framework's existing provider
registry + the `settings` config store (no duplication).

## API
- `providers.Active(k, capability, default) string` — resolve the selected backend.
- `providers.Use(k, capability, name, impl, def) bool` — a backend installs itself iff selected.
- `providers.Value(k, capability, name, key, default, secret) string` — resolve a config value (secrets env-only).

## Selection order
`env (TOGO_<CAP>_PROVIDER)` → `settings (providers.<cap>.active)` → `default`.
