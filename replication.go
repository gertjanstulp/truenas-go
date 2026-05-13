package truenas

import "encoding/json"

// import "encoding/json"

// ReplicationTaskResponse represents a replication task from the API.
type ReplicationTaskResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Direction   string `json:"direction"`
	Description string `json:"description"`
	Transport   string

	Attributes json.RawMessage `json:"attributes"` // Can be object or false
}
