package providers

// Capability names — the shared vocabulary for provider selection. Some already
// have a contract in togo core (contracts.go); the selection helper in this
// package works for ANY capability string, existing or new.
const (
	// --- already defined by togo core (contracts.go) — reuse, don't redefine ---
	CapCache   = "cache"    // togo.Cache
	CapQueue   = "queue"    // togo.Queue      → queue-nats (default), queue-rabbitmq, queue-kafka
	CapStorage = "storage"  // togo.Storage    → storage-s3, storage-r2, storage-gdrive
	CapRealtime = "realtime" // togo.Broker    → realtime-nats, realtime-grpc

	// --- new capabilities (contract lives in the owning plugin) ---
	CapImplement = "impl"    // autopilot: issue → code. claude (default), omnigent
	CapExecute   = "exec"    // autopilot: isolated env. local (default), coder
	CapCompute   = "compute" // worker: batch/stream jobs. worker-local (default), beam/spark/flink/databricks
	CapData      = "data"    // db: governed read/write. pg (default), iceberg/bigquery/databricks
	CapCatalog   = "catalog" // catalog: Unity Catalog / Iceberg REST
	CapShare     = "share"   // sharing: OpenSharing / Delta Sharing
)
