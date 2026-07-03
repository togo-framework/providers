// Package providers adds ONE thing on top of togo's existing provider system:
// runtime, config-driven SELECTION of which backend is active for a capability —
// swappable over CLI / .env / settings without recompiling.
//
// It deliberately reuses what the framework already has and invents nothing new:
//
//   - the registry            → togo.RegisterProviderFunc / applyProviders
//   - the capability contracts → togo.Cache / Queue / Storage / Broker / Translator
//     (new ones like Implement/Execute/Compute live in the plugin that owns them)
//   - the container           → k.Set(capability, impl) / k.Get(capability)
//   - the config store         → github.com/togo-framework/settings (scope "providers")
//   - the <cap>-<backend> convention → storage-s3, realtime-nats, queue-nats, …
//
// A backend plugin calls Use() in its provider func: it installs itself as the
// capability service only when it is the selected (or default) backend. Selection
// resolves highest-first: (1) env TOGO_<CAP>_PROVIDER, (2) settings
// providers.<cap>.active, (3) the caller's default.
package providers

import (
	"context"
	"os"
	"strings"

	"github.com/togo-framework/togo"
	"github.com/togo-framework/settings"
)

// Active returns the selected backend name for a capability, or def if none is
// configured. Resolution order: env → settings("providers") → def.
func Active(k *togo.Kernel, capability, def string) string {
	if v := os.Getenv("TOGO_" + up(capability) + "_PROVIDER"); v != "" {
		return v
	}
	if s := settings.FromKernel(k); s != nil {
		var v string
		if ok, _ := s.Get(context.Background(), "providers", capability+".active", &v); ok && v != "" {
			return v
		}
	}
	return def
}

// Use installs impl as the capability's kernel service IF this backend is the
// selected one — or, when def is true, if nothing else is selected. Returns
// whether it took effect. Call it from a backend plugin's provider func:
//
//	togo.RegisterProviderFunc("db-iceberg", togo.PriorityService, func(k *togo.Kernel) error {
//	    providers.Use(k, providers.CapData, "iceberg", newIceberg(k), false)
//	    return nil
//	})
func Use(k *togo.Kernel, capability, name string, impl any, def bool) bool {
	active := Active(k, capability, "")
	if active == name || (active == "" && def) {
		k.Set(capability, impl)
		if k.Log != nil {
			k.Log.Info("provider selected", "capability", capability, "backend", name)
		}
		return true
	}
	return false
}

// Value resolves a provider-scoped config value: env (VENDOR_KEY / TOGO_<CAP>_<KEY>)
// → settings("providers", "<cap>.<name>.<key>") → def. Secrets never read from
// settings (they resolve from env/.env only), so togo.yaml/DB never hold them.
func Value(k *togo.Kernel, capability, name, key, def string, secret bool) string {
	// Try vendor-native env first (CODER_URL, OMNIGENT_API), then the namespaced
	// form (TOGO_<CAP>_<KEY>) so both conventions work.
	for _, env := range []string{up(name) + "_" + up(key), "TOGO_" + up(capability) + "_" + up(key)} {
		if v := os.Getenv(env); v != "" {
			return v
		}
	}
	if !secret {
		if s := settings.FromKernel(k); s != nil {
			var v string
			if ok, _ := s.Get(context.Background(), "providers", capability+"."+name+"."+key, &v); ok && v != "" {
				return v
			}
		}
	}
	return def
}

func up(s string) string { return strings.ToUpper(strings.ReplaceAll(s, "-", "_")) }
