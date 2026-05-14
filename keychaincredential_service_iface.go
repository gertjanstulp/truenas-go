package truenas

import "context"

// KeychainCredentialServiceAPI defines the interface for cloud sync credential and task operations.
type KeychainCredentialServiceAPI interface {
	CreateSSHKeyPair(ctx context.Context, opts CreateSSHKeyPairOpts) (*SSHKeyPair, error)
	GetSSHKeyPair(ctx context.Context, id int64) (*SSHKeyPair, error)
	ListSSHKeyPairs(ctx context.Context) ([]SSHKeyPair, error)
	UpdateSSHKeyPair(ctx context.Context, id int64, opts UpdateSSHKeyPairOpts) (*SSHKeyPair, error)
	DeleteSSHKeyPair(ctx context.Context, id int64) error
}

// Compile-time checks.
var _ KeychainCredentialServiceAPI = (*KeychainCredentialService)(nil)
var _ KeychainCredentialServiceAPI = (*MockKeychainCredentialService)(nil)

// MockKeychainCredentialService is a test double for KeychainCredentialServiceAPI.
type MockKeychainCredentialService struct {
	CreateSSHKeyPairFunc func(ctx context.Context, opts CreateSSHKeyPairOpts) (*SSHKeyPair, error)
	GetSSHKeyPairFunc    func(ctx context.Context, id int64) (*SSHKeyPair, error)
	ListSSHKeyPairsFunc  func(ctx context.Context) ([]SSHKeyPair, error)
	UpdateSSHKeyPairFunc func(ctx context.Context, id int64, opts UpdateSSHKeyPairOpts) (*SSHKeyPair, error)
	DeleteSSHKeyPairFunc func(ctx context.Context, id int64) error
}

func (m *MockKeychainCredentialService) CreateSSHKeyPair(ctx context.Context, opts CreateSSHKeyPairOpts) (*SSHKeyPair, error) {
	if m.CreateSSHKeyPairFunc != nil {
		return m.CreateSSHKeyPairFunc(ctx, opts)
	}
	return nil, nil
}

func (m *MockKeychainCredentialService) GetSSHKeyPair(ctx context.Context, id int64) (*SSHKeyPair, error) {
	if m.GetSSHKeyPairFunc != nil {
		return m.GetSSHKeyPairFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockKeychainCredentialService) ListSSHKeyPairs(ctx context.Context) ([]SSHKeyPair, error) {
	if m.ListSSHKeyPairsFunc != nil {
		return m.ListSSHKeyPairsFunc(ctx)
	}
	return nil, nil
}

func (m *MockKeychainCredentialService) UpdateSSHKeyPair(ctx context.Context, id int64, opts UpdateSSHKeyPairOpts) (*SSHKeyPair, error) {
	if m.UpdateSSHKeyPairFunc != nil {
		return m.UpdateSSHKeyPairFunc(ctx, id, opts)
	}
	return nil, nil
}

func (m *MockKeychainCredentialService) DeleteSSHKeyPair(ctx context.Context, id int64) error {
	if m.DeleteSSHKeyPairFunc != nil {
		return m.DeleteSSHKeyPairFunc(ctx, id)
	}
	return nil
}
