package operator

// NodeCapacity represents the transformer or line limit at a specific grid node.
type NodeCapacity struct {
	NodeID     string  `json:"node_id"`
	CapacityKW float64 `json:"capacity_kw"`
}

// DeviceMapping links a SEP 2.0 device (by LFDI/mRID) to a specific grid node.
type DeviceMapping struct {
	DeviceLFDI string `json:"device_lfdi"`
	NodeID     string `json:"node_id"`
}

// FeederTopology represents the complete state of the feeder for dispatch decisions.
type FeederTopology struct {
	Capacities []NodeCapacity  `json:"capacities"`
	Mappings   []DeviceMapping `json:"mappings"`
}
