package truenas

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestKeychainCredentialService_CreateSSHKeyPair(t *testing.T) {
	callCount := 0
	testData := sampleSSHKeyPair()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "keychaincredential.create" {
						t.Errorf("expected method keychaincredential.create, got %s", method)
					}
					p := params.(map[string]any)
					if p["name"] != testData.Name {
						t.Errorf("expected name '%s', got '%s'", testData.Name, p["name"])
					}
					if p["type"] != KeychainCredentialTypeSSHKeyPair {
						t.Errorf("expected type '%s', got '%s'", KeychainCredentialTypeSSHKeyPair, p["type"])
					}
					if _, ok := p["attributes"]; !ok {
						t.Error("expected attributes in params")
					} else {
						attrs := p["attributes"].(map[string]interface{})
						if attrs["public_key"] != testData.PublicKey {
							t.Errorf("expected attrs.public_key '%s', got '%s'", testData.PublicKey, p["public_key"])
						}
						if attrs["private_key"] != testData.PrivateKey {
							t.Errorf("expected attrs.private_key '%s', got '%s'", testData.PrivateKey, p["private_key"])
						}
					}

					return json.RawMessage(`{"id": 1}`), nil
				}
				return sampleSSHKeyPairJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPair, err := svc.CreateSSHKeyPair(context.Background(), CreateSSHKeyPairOpts{
		Name:       testData.Name,
		PublicKey:  testData.PublicKey,
		PrivateKey: testData.PrivateKey,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshKeyPair == nil {
		t.Fatal("expected non-nil sshkeypair")
	}
	if sshKeyPair.ID != 1 {
		t.Errorf("expected ID %d, got %d", testData.ID, sshKeyPair.ID)
	}
	if sshKeyPair.Name != "ssh keypair name" {
		t.Errorf("expected name '%s', got '%s'", testData.Name, sshKeyPair.Name)
	}
	if sshKeyPair.PublicKey != "public key" {
		t.Errorf("expected publickey '%s', got '%s'", testData.PublicKey, sshKeyPair.PublicKey)
	}
	if sshKeyPair.PrivateKey != "private key" {
		t.Errorf("expected privatekey '%s', got '%s'", testData.PrivateKey, sshKeyPair.PrivateKey)
	}
}

func TestKeychainCredentialService_CreateSSHKeyPair_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("connection refused")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPair, err := svc.CreateSSHKeyPair(context.Background(), CreateSSHKeyPairOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
	if sshKeyPair != nil {
		t.Error("expected nil sshkeypair on error")
	}
}

func TestKeychainCredentialService_CreateSSHKeyPair_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.CreateSSHKeyPair(context.Background(), CreateSSHKeyPairOpts{})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_GetSSHKeyPair(t *testing.T) {
	testData := sampleSSHKeyPair()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.query" {
					t.Errorf("expected method keychaincredential.query, got %s", method)
				}
				return sampleSSHKeyPairJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPair, err := svc.GetSSHKeyPair(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshKeyPair == nil {
		t.Fatal("expected non-nil sshKeyPair")
	}
	if sshKeyPair.ID != testData.ID {
		t.Errorf("expected ID %d, got %d", testData.ID, sshKeyPair.ID)
	}
	if sshKeyPair.Name != testData.Name {
		t.Errorf("expected name '%s', got '%s'", testData.Name, sshKeyPair.Name)
	}
	if sshKeyPair.PublicKey != testData.PublicKey {
		t.Errorf("expected publickey '%s', got '%s'", testData.PublicKey, sshKeyPair.PublicKey)
	}
	if sshKeyPair.PrivateKey != testData.PrivateKey {
		t.Errorf("expected privatekey '%s', got '%s'", testData.PrivateKey, sshKeyPair.PrivateKey)
	}
}

func TestKeychainCredentialService_GetSSHKeyPair_NotFound(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPair, err := svc.GetSSHKeyPair(context.Background(), 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshKeyPair != nil {
		t.Error("expected nil sshkeypair for not found")
	}
}

func TestKeychainCredentialService_GetSSHKeyPair_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("timeout")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.GetSSHKeyPair(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_GetSSHKeyPair_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.GetSSHKeyPair(context.Background(), 1)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_ListSSHKeyPairs(t *testing.T) {
	testData := sampleSSHKeyPairs()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.query" {
					t.Errorf("expected method keychaincredential.query, got %s", method)
				}
				if params == nil {
					t.Error("expected filter params for ListSSHKeyPairs")
				}
				return json.RawMessage(sampleSSHKeyPairsJSON()), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPairs, err := svc.ListSSHKeyPairs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sshKeyPairs) != 2 {
		t.Fatalf("expected 2 sshKeyPairs, got %d", len(sshKeyPairs))
	}

	if sshKeyPairs[0].ID != testData[0].ID {
		t.Errorf("expected ID[0] %d, got %d", testData[0].ID, sshKeyPairs[0].ID)
	}
	if sshKeyPairs[0].Name != testData[0].Name {
		t.Errorf("expected name[0] '%s', got '%s'", testData[0].Name, sshKeyPairs[0].Name)
	}
	if sshKeyPairs[0].PublicKey != testData[0].PublicKey {
		t.Errorf("expected publickey[0] '%s', got '%s'", testData[0].PublicKey, sshKeyPairs[0].PublicKey)
	}
	if sshKeyPairs[0].PrivateKey != testData[0].PrivateKey {
		t.Errorf("expected privatekey[0] '%s', got '%s'", testData[0].PrivateKey, sshKeyPairs[0].PrivateKey)
	}

	if sshKeyPairs[1].ID != testData[1].ID {
		t.Errorf("expected ID[1] %d, got %d", testData[1].ID, sshKeyPairs[1].ID)
	}
	if sshKeyPairs[1].Name != testData[1].Name {
		t.Errorf("expected name[1] '%s', got '%s'", testData[1].Name, sshKeyPairs[1].Name)
	}
	if sshKeyPairs[1].PublicKey != testData[1].PublicKey {
		t.Errorf("expected publickey[1] '%s', got '%s'", testData[1].PublicKey, sshKeyPairs[1].PublicKey)
	}
	if sshKeyPairs[1].PrivateKey != testData[1].PrivateKey {
		t.Errorf("expected privatekey[1] '%s', got '%s'", testData[1].PrivateKey, sshKeyPairs[1].PrivateKey)
	}
}

func TestKeychainCredentialService_ListSSHKeyPairs_Empty(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPairs, err := svc.ListSSHKeyPairs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sshKeyPairs) != 0 {
		t.Errorf("expected 0 sshkeypairs, got %d", len(sshKeyPairs))
	}
}

func TestKeychainCredentialService_ListSSHKeyPairs_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("network error")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.ListSSHKeyPairs(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_ListSSHKeyPairs_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.ListSSHKeyPairs(context.Background())
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_UpdateSSHKeyPair(t *testing.T) {
	callCount := 0
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "keychaincredential.update" {
						t.Errorf("expected method keychaincredential.update, got %s", method)
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
					p := slice[1].(map[string]any)
					if p["name"] != "Updated keypair" {
						t.Errorf("expected name 'Updated keypair', got %v", p["name"])
					}
					return json.RawMessage(`{"id": 1}`), nil
				}
				return sampleSSHKeyPairJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshKeyPair, err := svc.UpdateSSHKeyPair(context.Background(), 1, UpdateSSHKeyPairOpts{
		Name: "Updated keypair",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshKeyPair == nil {
		t.Fatal("expected non-nil sshkeypair")
	}
	if sshKeyPair.ID != 1 {
		t.Errorf("expected ID 1, got %d", sshKeyPair.ID)
	}
}

func TestKeychainCredentialService_UpdateSSHKeyPair_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("not found")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.UpdateSSHKeyPair(context.Background(), 999, UpdateSSHKeyPairOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_DeleteSSHKeyPair(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.delete" {
					t.Errorf("expected method keychaincredential.delete, got %s", method)
				}
				id, ok := params.(int64)
				if !ok || id != 5 {
					t.Errorf("expected id 5, got %v", params)
				}
				return nil, nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	err := svc.DeleteSSHKeyPair(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeychainCredentialService_DeleteSSHKeyPair_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("permission denied")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	err := svc.DeleteSSHKeyPair(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_CreateSSHCredential(t *testing.T) {
	callCount := 0
	testData := sampleSSHCredential()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "keychaincredential.create" {
						t.Errorf("expected method keychaincredential.create, got %s", method)
					}
					p := params.(map[string]any)
					if p["name"] != testData.Name {
						t.Errorf("expected name '%s', got '%s'", testData.Name, p["name"])
					}
					if p["type"] != KeychainCredentialTypeSSHCredential {
						t.Errorf("expected type '%s', got '%s'", KeychainCredentialTypeSSHCredential, p["type"])
					}

					if _, ok := p["attributes"]; !ok {
						t.Error("expected attributes in params")
					} else {
						attrs := p["attributes"].(map[string]interface{})
						if attrs["host"] != testData.Host {
							t.Errorf("expected host '%s', got '%s'", testData.Host, p["host"])
						}
						if attrs["port"] != testData.Port {
							t.Errorf("expected port '%d', got '%d'", testData.Port, p["port"])
						}
						if attrs["username"] != testData.Username {
							t.Errorf("expected username '%s', got '%s'", testData.Username, p["username"])
						}
						if attrs["remote_host_key"] != testData.RemoteHostKey {
							t.Errorf("expected remote_host_key '%s', got '%s'", testData.RemoteHostKey, p["remote_host_key"])
						}
						if attrs["connect_timeout"] != testData.ConnectTimeout {
							t.Errorf("expected connect_timeout '%d', got '%d'", testData.ConnectTimeout, p["connect_timeout"])
						}
						if attrs["private_key"] != testData.SSHKeyPairID {
							t.Errorf("expected private_key '%d', got '%d'", testData.SSHKeyPairID, p["private_key"])
						}
					}

					return json.RawMessage(`{"id": 1}`), nil
				}
				return sampleSSHCredentialJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredential, err := svc.CreateSSHCredential(context.Background(), CreateSSHCredentialOpts{
		Name:           testData.Name,
		Host:           testData.Host,
		Port:           testData.Port,
		Username:       testData.Username,
		RemoteHostKey:  testData.RemoteHostKey,
		ConnectTimeout: testData.ConnectTimeout,
		SSHKeyPairID:   testData.SSHKeyPairID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshCredential == nil {
		t.Fatal("expected non-nil sshcredential")
	}
	if sshCredential.ID != 1 {
		t.Errorf("expected ID %d, got %d", testData.ID, sshCredential.ID)
	}
	if sshCredential.Name != testData.Name {
		t.Errorf("expected name '%s', got '%s'", testData.Name, sshCredential.Name)
	}
	if sshCredential.Host != testData.Host {
		t.Errorf("expected host '%s', got %s", testData.Host, sshCredential.Host)
	}
	if sshCredential.Port != testData.Port {
		t.Errorf("expected port '%d', got %d", testData.Port, sshCredential.Port)
	}
	if sshCredential.Username != testData.Username {
		t.Errorf("expected username '%s', got %s", testData.Username, sshCredential.Username)
	}
	if sshCredential.RemoteHostKey != testData.RemoteHostKey {
		t.Errorf("expected remote_host_key '%s', got %s", testData.RemoteHostKey, sshCredential.RemoteHostKey)
	}
	if sshCredential.ConnectTimeout != testData.ConnectTimeout {
		t.Errorf("expected connect_timeout '%d', got %d", testData.ConnectTimeout, sshCredential.ConnectTimeout)
	}
	if sshCredential.SSHKeyPairID != testData.SSHKeyPairID {
		t.Errorf("expected private_key '%d', got %d", testData.SSHKeyPairID, sshCredential.SSHKeyPairID)
	}
}

func TestKeychainCredentialService_CreateSSHCredential_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("connection refused")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredential, err := svc.CreateSSHCredential(context.Background(), CreateSSHCredentialOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
	if sshCredential != nil {
		t.Error("expected nil sshkeypair on error")
	}
}

func TestKeychainCredentialService_CreateSSHCredential_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.CreateSSHCredential(context.Background(), CreateSSHCredentialOpts{})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_GetSSHCredential(t *testing.T) {
	testData := sampleSSHCredential()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.query" {
					t.Errorf("expected method keychaincredential.query, got %s", method)
				}
				return sampleSSHCredentialJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredential, err := svc.GetSSHCredential(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshCredential == nil {
		t.Fatal("expected non-nil sshCredential")
	}
	if sshCredential.ID != testData.ID {
		t.Errorf("expected ID %d, got %d", testData.ID, sshCredential.ID)
	}
	if sshCredential.Name != testData.Name {
		t.Errorf("expected name '%s', got '%s'", testData.Name, sshCredential.Name)
	}
	if sshCredential.Host != testData.Host {
		t.Errorf("expected host '%s', got %s", testData.Host, sshCredential.Host)
	}
	if sshCredential.Port != testData.Port {
		t.Errorf("expected port '%d', got %d", testData.Port, sshCredential.Port)
	}
	if sshCredential.Username != testData.Username {
		t.Errorf("expected username '%s', got %s", testData.Username, sshCredential.Username)
	}
	if sshCredential.RemoteHostKey != testData.RemoteHostKey {
		t.Errorf("expected remote_host_key '%s', got %s", testData.RemoteHostKey, sshCredential.RemoteHostKey)
	}
	if sshCredential.ConnectTimeout != testData.ConnectTimeout {
		t.Errorf("expected connect_timeout '%d', got %d", testData.ConnectTimeout, sshCredential.ConnectTimeout)
	}
	if sshCredential.SSHKeyPairID != testData.SSHKeyPairID {
		t.Errorf("expected private_key '%d', got %d", testData.SSHKeyPairID, sshCredential.SSHKeyPairID)
	}
}

func TestKeychainCredentialService_GetSSHCredential_NotFound(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredential, err := svc.GetSSHCredential(context.Background(), 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshCredential != nil {
		t.Error("expected nil sshkeypair for not found")
	}
}

func TestKeychainCredentialService_GetSSHCredential_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("timeout")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.GetSSHCredential(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_GetSSHCredential_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.GetSSHCredential(context.Background(), 1)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_ListSSHCredentials(t *testing.T) {
	testData := sampleSSHCredentials()
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.query" {
					t.Errorf("expected method keychaincredential.query, got %s", method)
				}
				if params == nil {
					t.Error("expected filter params for ListSSHCredentials")
				}
				return json.RawMessage(sampleSSHCredentialsJSON()), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredentials, err := svc.ListSSHCredentials(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sshCredentials) != 2 {
		t.Fatalf("expected 2 sshCredentials, got %d", len(sshCredentials))
	}

	if sshCredentials[0].ID != testData[0].ID {
		t.Errorf("expected ID[0] %d, got %d", testData[0].ID, sshCredentials[0].ID)
	}
	if sshCredentials[0].Name != testData[0].Name {
		t.Errorf("expected name[0] '%s', got '%s'", testData[0].Name, sshCredentials[0].Name)
	}
	if sshCredentials[0].Host != testData[0].Host {
		t.Errorf("expected host '%s', got %s", testData[0].Host, sshCredentials[0].Host)
	}
	if sshCredentials[0].Port != testData[0].Port {
		t.Errorf("expected port '%d', got %d", testData[0].Port, sshCredentials[0].Port)
	}
	if sshCredentials[0].Username != testData[0].Username {
		t.Errorf("expected username '%s', got %s", testData[0].Username, sshCredentials[0].Username)
	}
	if sshCredentials[0].RemoteHostKey != testData[0].RemoteHostKey {
		t.Errorf("expected remote_host_key '%s', got %s", testData[0].RemoteHostKey, sshCredentials[0].RemoteHostKey)
	}
	if sshCredentials[0].ConnectTimeout != testData[0].ConnectTimeout {
		t.Errorf("expected connect_timeout '%d', got %d", testData[0].ConnectTimeout, sshCredentials[0].ConnectTimeout)
	}
	if sshCredentials[0].SSHKeyPairID != testData[0].SSHKeyPairID {
		t.Errorf("expected private_key '%d', got %d", testData[0].SSHKeyPairID, sshCredentials[0].SSHKeyPairID)
	}

	if sshCredentials[1].ID != testData[1].ID {
		t.Errorf("expected ID[1] %d, got %d", testData[1].ID, sshCredentials[1].ID)
	}
	if sshCredentials[1].Name != testData[1].Name {
		t.Errorf("expected name[1] '%s', got '%s'", testData[1].Name, sshCredentials[1].Name)
	}
	if sshCredentials[1].Host != testData[1].Host {
		t.Errorf("expected host '%s', got %s", testData[1].Host, sshCredentials[1].Host)
	}
	if sshCredentials[1].Port != testData[1].Port {
		t.Errorf("expected port '%d', got %d", testData[1].Port, sshCredentials[1].Port)
	}
	if sshCredentials[1].Username != testData[1].Username {
		t.Errorf("expected username '%s', got %s", testData[1].Username, sshCredentials[1].Username)
	}
	if sshCredentials[1].RemoteHostKey != testData[1].RemoteHostKey {
		t.Errorf("expected remote_host_key '%s', got %s", testData[1].RemoteHostKey, sshCredentials[1].RemoteHostKey)
	}
	if sshCredentials[1].ConnectTimeout != testData[1].ConnectTimeout {
		t.Errorf("expected connect_timeout '%d', got %d", testData[1].ConnectTimeout, sshCredentials[1].ConnectTimeout)
	}
	if sshCredentials[1].SSHKeyPairID != testData[1].SSHKeyPairID {
		t.Errorf("expected private_key '%d', got %d", testData[1].SSHKeyPairID, sshCredentials[1].SSHKeyPairID)
	}
}

func TestKeychainCredentialService_ListSSHCredentials_Empty(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`[]`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredentials, err := svc.ListSSHCredentials(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sshCredentials) != 0 {
		t.Errorf("expected 0 sshkeypairs, got %d", len(sshCredentials))
	}
}

func TestKeychainCredentialService_ListSSHCredentials_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("network error")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.ListSSHCredentials(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_ListSSHCredentials_ParseError(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return json.RawMessage(`not json`), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.ListSSHCredentials(context.Background())
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestKeychainCredentialService_UpdateSSHCredential(t *testing.T) {
	callCount := 0
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				callCount++
				if callCount == 1 {
					if method != "keychaincredential.update" {
						t.Errorf("expected method keychaincredential.update, got %s", method)
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
					p := slice[1].(map[string]any)
					if p["name"] != "Updated keypair" {
						t.Errorf("expected name 'Updated keypair', got %v", p["name"])
					}
					return json.RawMessage(`{"id": 1}`), nil
				}
				return sampleSSHCredentialJSON(), nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	sshCredential, err := svc.UpdateSSHCredential(context.Background(), 1, UpdateSSHCredentialOpts{
		Name: "Updated keypair",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sshCredential == nil {
		t.Fatal("expected non-nil sshkeypair")
	}
	if sshCredential.ID != 1 {
		t.Errorf("expected ID 1, got %d", sshCredential.ID)
	}
}

func TestKeychainCredentialService_UpdateSSHCredential_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("not found")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	_, err := svc.UpdateSSHCredential(context.Background(), 999, UpdateSSHCredentialOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKeychainCredentialService_DeleteSSHCredential(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				if method != "keychaincredential.delete" {
					t.Errorf("expected method keychaincredential.delete, got %s", method)
				}
				id, ok := params.(int64)
				if !ok || id != 5 {
					t.Errorf("expected id 5, got %v", params)
				}
				return nil, nil
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	err := svc.DeleteSSHCredential(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeychainCredentialService_DeleteSSHCredential_Error(t *testing.T) {
	mock := &mockAsyncCaller{
		mockCaller: mockCaller{
			callFunc: func(ctx context.Context, method string, params any) (json.RawMessage, error) {
				return nil, errors.New("permission denied")
			},
		},
	}

	svc := NewKeychainCredentialService(mock, Version{Major: 25, Minor: 4})
	err := svc.DeleteSSHCredential(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
