package truenas

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// sampleReplicationJSON returns a JSON array response for a replication query.
func sampleReplicationJSON() json.RawMessage {
	return json.RawMessage(`[{
		"id": 1,
		"description": "Daily replication"
	}]`)
}

// sampleReplicationSingleJSON returns a single JSON object response for get_instance.
// func sampleReplicationSingleJSON() json.RawMessage {
// 	return json.RawMessage(`{
// 		"id": 1,
// 		"description": "Daily replication"
// 	}`)
// }

func TestReplicationService_Create(t *testing.T) {
	callCount := 0
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "replication.create" {
						t.Errorf("expected method replication.create, got %s", method)
					}
					return json.RawMessage(`{"id": 1}`), nil
				}
				// Re-read via get_instance returns single object
				return sampleReplicationJSON(), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	replication, err := svc.CreateTask(context.Background(), CreateReplicationTaskOpts{
		Description: "Daily replication",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if replication == nil {
		t.Fatal("expected non-nil replication")
	}
	if replication.ID != 1 {
		t.Errorf("expected ID 1, got %d", replication.ID)
	}
}

func TestReplicationService_Create_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("connection refused")
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	replication, err := svc.CreateTask(context.Background(), CreateReplicationTaskOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
	if replication != nil {
		t.Error("expected nil replication on error")
	}
	if err.Error() != "connection refused" {
		t.Errorf("expected 'connection refused', got %q", err.Error())
	}
}

func TestReplicationService_Create_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	_, err := svc.CreateTask(context.Background(), CreateReplicationTaskOpts{})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestReplicationService_Get(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "replication.query" {
					t.Errorf("expected method replication.query, got %s", method)
				}
				return sampleReplicationJSON(), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	replication, err := svc.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if replication == nil {
		t.Fatal("expected non-nil replication")
	}
	if replication.ID != 1 {
		t.Errorf("expected ID 1, got %d", replication.ID)
	}
	if replication.Description != "Daily replication" {
		t.Errorf("expected description 'Daily replication', got %q", replication.Description)
	}
}

func TestReplicationService_Get_NotFound(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	replication, err := svc.GetTask(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected nil error for not-found, got: %v", err)
	}
	if replication != nil {
		t.Error("expected nil replication for not found")
	}
}

func TestReplicationService_Get_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("timeout")
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	_, err := svc.GetTask(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReplicationService_List(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "replication.query" {
					t.Errorf("expected method replication.query, got %s", method)
				}
				if params != nil {
					t.Error("expected nil params for List")
				}
				return json.RawMessage(`[
				{"id": 1},
				{"id": 2}
			]`), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	jobs, err := svc.ListTasks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
	if jobs[0].ID != 1 {
		t.Errorf("expected first replication ID 1, got %d", jobs[0].ID)
	}
}

func TestReplicationService_List_Empty(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	jobs, err := svc.ListTasks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestReplicationService_List_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("network error")
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	_, err := svc.ListTasks(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReplicationService_Update(t *testing.T) {
	callCount := 0
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "replication.update" {
						t.Errorf("expected method replication.update, got %s", method)
					}
					slice, ok := params.([]any)
					if !ok {
						t.Fatal("expected []any params for update")
					}
					if len(slice) != 2 {
						t.Fatalf("expected 2 elements, got %d", len(slice))
					}
					id, ok := slice[0].(int64)
					if !ok || id != 1 {
						t.Errorf("expected id 1, got %v", slice[0])
					}
					return json.RawMessage(`{"id": 1}`), nil
				}
				// Re-read via get_instance returns single object
				return sampleReplicationJSON(), nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	replication, err := svc.UpdateTask(context.Background(), 1, UpdateReplicationTaskOpts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if replication == nil {
		t.Fatal("expected non-nil replication")
	}
	if replication.ID != 1 {
		t.Errorf("expected ID 1, got %d", replication.ID)
	}
}

func TestReplicationService_Update_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("not found")
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	_, err := svc.UpdateTask(context.Background(), 999, UpdateReplicationTaskOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReplicationService_Delete(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "replication.delete" {
					t.Errorf("expected method replication.delete, got %s", method)
				}
				id, ok := params.(int64)
				if !ok || id != 5 {
					t.Errorf("expected id 5, got %v", params)
				}
				return nil, nil
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	err := svc.DeleteTask(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplicationService_Delete_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("permission denied")
			},
		},
	}

	svc := NewReplicationService(mock, Version{})
	err := svc.DeleteTask(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
