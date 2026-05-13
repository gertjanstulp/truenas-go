package truenas

import (
	"context"
	"encoding/json"
	"fmt"
)

// ReplicationTask is the user-facing representation of a replication task.
type ReplicationTask struct {
	ID          int64
	Name        string
	Direction   string
	Description string
	Attributes  map[string]any
}

// CreateReplicationTaskOpts contains options for creating a replication task.
type CreateReplicationTaskOpts struct {
	Description string
	Attributes  map[string]any
}

// UpdateReplicationTaskOpts contains options for updating a replication task.
type UpdateReplicationTaskOpts = CreateReplicationTaskOpts

// ReplicationService provides typed methods for the replication.* API namespace.
type ReplicationService struct {
	client  AsyncCaller
	version Version
}

// NewReplicationService creates a new ReplicationService.
func NewReplicationService(c AsyncCaller, v Version) *ReplicationService {
	return &ReplicationService{client: c, version: v}
}

// CreateReplicationTask creates a replication task and returns the full object.
func (s *ReplicationService) CreateTask(ctx context.Context, opts CreateReplicationTaskOpts) (*ReplicationTask, error) {
	params := replicationTaskOptsToParams(opts)

	result, err := s.client.Call(ctx, "replication.create", params)
	if err != nil {
		return nil, err
	}

	var createResp struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(result, &createResp); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}

	return s.GetTask(ctx, createResp.ID)
}

// GetTask returns a replication task by ID, or nil if not found.
func (s *ReplicationService) GetTask(ctx context.Context, id int64) (*ReplicationTask, error) {
	filter := [][]any{{"id", "=", id}}
	result, err := s.client.Call(ctx, "replication.query", filter)
	if err != nil {
		return nil, err
	}

	var tasks []ReplicationTaskResponse
	if err := json.Unmarshal(result, &tasks); err != nil {
		return nil, fmt.Errorf("parse query response: %w", err)
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	task := replicationTaskFromResponse(tasks[0])
	return &task, nil
}

// ListTasks returns all replication tasks.
func (s *ReplicationService) ListTasks(ctx context.Context) ([]ReplicationTask, error) {
	result, err := s.client.Call(ctx, "replication.query", nil)
	if err != nil {
		return nil, err
	}

	var responses []ReplicationTaskResponse
	if err := json.Unmarshal(result, &responses); err != nil {
		return nil, fmt.Errorf("parse query response: %w", err)
	}

	tasks := make([]ReplicationTask, len(responses))
	for i, resp := range responses {
		tasks[i] = replicationTaskFromResponse(resp)
	}
	return tasks, nil
}

// UpdateTask updates a replication task and returns the full object.
func (s *ReplicationService) UpdateTask(ctx context.Context, id int64, opts UpdateReplicationTaskOpts) (*ReplicationTask, error) {
	params := replicationTaskOptsToParams(opts)

	_, err := s.client.Call(ctx, "replication.update", []any{id, params})
	if err != nil {
		return nil, err
	}

	return s.GetTask(ctx, id)
}

// DeleteTask deletes a replication task by ID.
func (s *ReplicationService) DeleteTask(ctx context.Context, id int64) error {
	_, err := s.client.Call(ctx, "replication.delete", id)
	return err
}

// taskOptsToParams converts CreateReplicationTaskOpts to API parameters.
func replicationTaskOptsToParams(opts CreateReplicationTaskOpts) map[string]any {
	params := map[string]any{
		"description": opts.Description,
	}

	attrs := opts.Attributes
	if attrs == nil {
		attrs = map[string]any{}
	}
	params["attributes"] = attrs

	return params
}

// taskFromResponse converts a wire-format ReplicationTaskResponse to a user-facing ReplicationTask.
func replicationTaskFromResponse(resp ReplicationTaskResponse) ReplicationTask {
	task := ReplicationTask{
		ID:          resp.ID,
		Name:        resp.Name,
		Direction:   resp.Direction,
		Description: resp.Description,
	}

	// Handle attributes - can be false in API response, so ignore unmarshal errors
	if len(resp.Attributes) > 0 {
		var attrs map[string]any
		if err := json.Unmarshal(resp.Attributes, &attrs); err == nil {
			task.Attributes = attrs
		}
	}

	return task
}
