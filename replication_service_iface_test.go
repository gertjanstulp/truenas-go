package truenas

import (
	"context"
	"testing"
)

func TestMockReplicationService_ImplementsInterface(t *testing.T) {
	var _ ReplicationServiceAPI = (*ReplicationService)(nil)
	var _ ReplicationServiceAPI = (*MockReplicationService)(nil)
}

func TestMockReplicationService_DefaultsToNil(t *testing.T) {
	mock := &MockReplicationService{}
	ctx := context.Background()

	cred, err := mock.GetTask(ctx, 1)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if cred != nil {
		t.Fatalf("expected nil result, got: %v", cred)
	}
}

func TestMockReplicationService_CallsFunc(t *testing.T) {
	called := false
	mock := &MockReplicationService{
		GetTaskFunc: func(ctx context.Context, id int64) (*ReplicationTask, error) {
			called = true
			return &ReplicationTask{ID: id}, nil
		},
	}

	cred, err := mock.GetTask(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected GetCredGetTaskFuncentialFunc to be called")
	}
	if cred.ID != 42 {
		t.Fatalf("expected ID 42, got %d", cred.ID)
	}
}
