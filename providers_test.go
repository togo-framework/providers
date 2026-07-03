package providers

import (
	"testing"

	"github.com/togo-framework/togo"
)

func TestActiveResolvesEnvThenDefault(t *testing.T) {
	k := togo.New()

	// no env, no settings → default wins
	if got := Active(k, CapData, "pg"); got != "pg" {
		t.Fatalf("Active default = %q, want pg", got)
	}
	// env override wins
	t.Setenv("TOGO_DATA_PROVIDER", "iceberg")
	if got := Active(k, CapData, "pg"); got != "iceberg" {
		t.Fatalf("Active env = %q, want iceberg", got)
	}
}

func TestUseSelectsOnlyTheActiveBackend(t *testing.T) {
	k := togo.New()
	t.Setenv("TOGO_DATA_PROVIDER", "iceberg")

	// the non-selected default backend must NOT install itself
	if Use(k, CapData, "pg", "PG_IMPL", true) {
		t.Fatal("pg installed itself despite iceberg being selected")
	}
	// the selected backend installs itself
	if !Use(k, CapData, "iceberg", "ICEBERG_IMPL", false) {
		t.Fatal("iceberg did not install itself despite being selected")
	}
	if v, ok := k.Get(CapData); !ok || v != "ICEBERG_IMPL" {
		t.Fatalf("k.Get(data) = %v,%v; want ICEBERG_IMPL", v, ok)
	}
}

func TestUseDefaultWhenNothingSelected(t *testing.T) {
	k := togo.New()
	// nothing selected → the default backend takes the slot
	if !Use(k, CapCompute, "worker-local", "LOCAL", true) {
		t.Fatal("default backend did not take the empty slot")
	}
	if v, _ := k.Get(CapCompute); v != "LOCAL" {
		t.Fatalf("k.Get(compute) = %v; want LOCAL", v)
	}
}

func TestValueEnvResolution(t *testing.T) {
	k := togo.New()
	t.Setenv("CODER_URL", "https://coder.example.com")
	if got := Value(k, CapExecute, "coder", "url", "http://localhost", false); got != "https://coder.example.com" {
		t.Fatalf("Value env = %q", got)
	}
	// secret falls back to default when env unset (never reads settings)
	if got := Value(k, CapExecute, "coder", "token", "", true); got != "" {
		t.Fatalf("Value secret default = %q, want empty", got)
	}
}
