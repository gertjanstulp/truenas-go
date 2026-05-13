package truenas

import "context"

// ReplicationServiceAPI defines the interface for replication credential and task operations.
type ReplicationServiceAPI interface {
	CreateTask(ctx context.Context, opts CreateReplicationTaskOpts) (*ReplicationTask, error)
	GetTask(ctx context.Context, id int64) (*ReplicationTask, error)
	ListTasks(ctx context.Context) ([]ReplicationTask, error)
	UpdateTask(ctx context.Context, id int64, opts UpdateReplicationTaskOpts) (*ReplicationTask, error)
	DeleteTask(ctx context.Context, id int64) error
}

// Compile-time checks.
var _ ReplicationServiceAPI = (*ReplicationService)(nil)
var _ ReplicationServiceAPI = (*MockReplicationService)(nil)

// MockReplicationService is a test double for ReplicationServiceAPI.
type MockReplicationService struct {
	CreateTaskFunc func(ctx context.Context, opts CreateReplicationTaskOpts) (*ReplicationTask, error)
	GetTaskFunc    func(ctx context.Context, id int64) (*ReplicationTask, error)
	ListTasksFunc  func(ctx context.Context) ([]ReplicationTask, error)
	UpdateTaskFunc func(ctx context.Context, id int64, opts UpdateReplicationTaskOpts) (*ReplicationTask, error)
	DeleteTaskFunc func(ctx context.Context, id int64) error
}

func (m *MockReplicationService) CreateTask(ctx context.Context, opts CreateReplicationTaskOpts) (*ReplicationTask, error) {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(ctx, opts)
	}
	return nil, nil
}

func (m *MockReplicationService) GetTask(ctx context.Context, id int64) (*ReplicationTask, error) {
	if m.GetTaskFunc != nil {
		return m.GetTaskFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockReplicationService) ListTasks(ctx context.Context) ([]ReplicationTask, error) {
	if m.ListTasksFunc != nil {
		return m.ListTasksFunc(ctx)
	}
	return nil, nil
}

func (m *MockReplicationService) UpdateTask(ctx context.Context, id int64, opts UpdateReplicationTaskOpts) (*ReplicationTask, error) {
	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(ctx, id, opts)
	}
	return nil, nil
}

func (m *MockReplicationService) DeleteTask(ctx context.Context, id int64) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}
	return nil
}
