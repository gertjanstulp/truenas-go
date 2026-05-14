package truenas

import (
	"context"
	"encoding/json"
	"fmt"
)

// SSHKeyPair is the user-facing representation of a TrueNAS ssh keypair.
type SSHKeyPair struct {
	ID         int64
	Name       string
	PrivateKey string
	PublicKey  string
}

// CreateSSHKeyPairOpts contains options for creating a ssh keypair.
type CreateSSHKeyPairOpts struct {
	Name       string
	PublicKey  string
	PrivateKey string
}

// UpdateSSHKeyPairOpts contains options for updating a ssh keypair.
type UpdateSSHKeyPairOpts = CreateSSHKeyPairOpts

// SSHCredential is the user-facing representation of a TrueNAS ssh credential.
type SSHCredential struct {
	ID             int64
	Name           string
	Host           string
	Port           int32
	Username       string
	RemoteHostKey  string
	ConnectTimeout int
	SSHKeyPairID   int64
}

// CreateSSHCredentialOpts contains options for creating an ssh keypair.
type CreateSSHCredentialOpts struct {
	Name           string
	Host           string
	Port           int32
	Username       string
	RemoteHostKey  string
	ConnectTimeout int
	SSHKeyPairID   int64
}

// UpdateSSHCredentialOpts contains options for updating an ssh credential.
type UpdateSSHCredentialOpts = CreateSSHCredentialOpts

// KeychainCredentialService provides typed methods for the keychaincredential.query API namespaces.
type KeychainCredentialService struct {
	client  Caller
	version Version
}

// NewKeychainCredentialService creates a new KeychainCredentialService.
func NewKeychainCredentialService(c Caller, v Version) *KeychainCredentialService {
	return &KeychainCredentialService{client: c, version: v}
}

// CreateSSHKeyPair creates an ssh keypair and returns the full object.
func (s *KeychainCredentialService) CreateSSHKeyPair(ctx context.Context, opts CreateSSHKeyPairOpts) (*SSHKeyPair, error) {
	params := sshKeyPairOptsToParams(opts, true)

	result, err := s.client.Call(ctx, "keychaincredential.create", params)
	if err != nil {
		return nil, err
	}

	var createResp struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(result, &createResp); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}

	return s.GetSSHKeyPair(ctx, createResp.ID)
}

// GetSSHKeyPair returns a SSH keypair by ID, or nil if not found.
func (s *KeychainCredentialService) GetSSHKeyPair(ctx context.Context, id int64) (*SSHKeyPair, error) {
	filter := [][]any{{"id", "=", id}}
	result, err := s.client.Call(ctx, "keychaincredential.query", filter)
	if err != nil {
		return nil, err
	}

	keyPairs, err := ParseSSHKeyPairs(result)
	if err != nil {
		return nil, err
	}

	if len(keyPairs) == 0 {
		return nil, nil
	}

	keyPair := sshKeyPairFromResponse(keyPairs[0])
	return &keyPair, nil
}

// ListSSHKeyPairs returns all keychain credentials.
func (s *KeychainCredentialService) ListSSHKeyPairs(ctx context.Context) ([]SSHKeyPair, error) {
	filter := [][]any{{"type", "=", KeychainCredentialTypeSSHKeyPair}}
	result, err := s.client.Call(ctx, "keychaincredential.query", filter)
	if err != nil {
		return nil, err
	}

	keyPairs, err := ParseSSHKeyPairs(result)
	if err != nil {
		return nil, err
	}

	sshKeyPairs := make([]SSHKeyPair, len(keyPairs))
	for i, resp := range keyPairs {
		sshKeyPairs[i] = sshKeyPairFromResponse(resp)
	}
	return sshKeyPairs, nil
}

// UpdateSSHKeyPair updates a SSH keypair and returns the full object.
func (s *KeychainCredentialService) UpdateSSHKeyPair(ctx context.Context, id int64, opts UpdateSSHKeyPairOpts) (*SSHKeyPair, error) {
	params := sshKeyPairOptsToParams(opts, false)

	_, err := s.client.Call(ctx, "keychaincredential.update", []any{id, params})
	if err != nil {
		return nil, err
	}

	return s.GetSSHKeyPair(ctx, id)
}

// DeleteSSHKeyPair deletes a SSH keypair by ID.
func (s *KeychainCredentialService) DeleteSSHKeyPair(ctx context.Context, id int64) error {
	_, err := s.client.Call(ctx, "keychaincredential.delete", id)
	return err
}

// CreateSSHCredential creates an ssh credential and returns the full object.
func (s *KeychainCredentialService) CreateSSHCredential(ctx context.Context, opts CreateSSHCredentialOpts) (*SSHCredential, error) {
	params := sshCredentialOptsToParams(opts, true)

	result, err := s.client.Call(ctx, "keychaincredential.create", params)
	if err != nil {
		return nil, err
	}

	var createResp struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(result, &createResp); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}

	return s.GetSSHCredential(ctx, createResp.ID)
}

// GetSSHCredential returns a SSH credential by ID, or nil if not found.
func (s *KeychainCredentialService) GetSSHCredential(ctx context.Context, id int64) (*SSHCredential, error) {
	filter := [][]any{{"id", "=", id}}
	result, err := s.client.Call(ctx, "keychaincredential.query", filter)
	if err != nil {
		return nil, err
	}

	credentials, err := ParseSSHCredentials(result)
	if err != nil {
		return nil, err
	}

	if len(credentials) == 0 {
		return nil, nil
	}

	credential := sshCredentialFromResponse(credentials[0])
	return &credential, nil
}

// ListSSHCredentials returns all keychain credentials.
func (s *KeychainCredentialService) ListSSHCredentials(ctx context.Context) ([]SSHCredential, error) {
	filter := [][]any{{"type", "=", KeychainCredentialTypeSSHCredential}}
	result, err := s.client.Call(ctx, "keychaincredential.query", filter)
	if err != nil {
		return nil, err
	}

	credentials, err := ParseSSHCredentials(result)
	if err != nil {
		return nil, err
	}

	sshCredentials := make([]SSHCredential, len(credentials))
	for i, resp := range credentials {
		sshCredentials[i] = sshCredentialFromResponse(resp)
	}
	return sshCredentials, nil
}

// UpdateSSHCredential updates a SSH credential and returns the full object.
func (s *KeychainCredentialService) UpdateSSHCredential(ctx context.Context, id int64, opts UpdateSSHCredentialOpts) (*SSHCredential, error) {
	params := sshCredentialOptsToParams(opts, false)

	_, err := s.client.Call(ctx, "keychaincredential.update", []any{id, params})
	if err != nil {
		return nil, err
	}

	return s.GetSSHCredential(ctx, id)
}

// DeleteSSHCredential deletes a SSH credential by ID.
func (s *KeychainCredentialService) DeleteSSHCredential(ctx context.Context, id int64) error {
	_, err := s.client.Call(ctx, "keychaincredential.delete", id)
	return err
}

// sshKeyPairOptsToParams converts CreateSSHKeyPairOpts to API parameters.
func sshKeyPairOptsToParams(opts CreateSSHKeyPairOpts, includeType bool) map[string]any {
	attrs := map[string]any{
		"public_key":  opts.PublicKey,
		"private_key": opts.PrivateKey,
	}
	params := map[string]any{
		"name":       opts.Name,
		"attributes": attrs,
	}

	if includeType {
		params["type"] = KeychainCredentialTypeSSHKeyPair
	}

	return params
}

// sshKeyPairFromResponse converts a wire-format SSHKeyPairResponse to a user-facing SSHKeyPair.
func sshKeyPairFromResponse(resp SSHKeyPairResponse) SSHKeyPair {
	return SSHKeyPair{
		ID:         resp.ID,
		Name:       resp.Name,
		PrivateKey: resp.PrivateKey,
		PublicKey:  resp.PublicKey,
	}
}

// sshCredentialOptsToParams converts CreateSSHCredentialOpts to API parameters.
func sshCredentialOptsToParams(opts CreateSSHCredentialOpts, includeType bool) map[string]any {
	attrs := map[string]any{
		"host":            opts.Host,
		"port":            opts.Port,
		"username":        opts.Username,
		"remote_host_key": opts.RemoteHostKey,
		"connect_timeout": opts.ConnectTimeout,
		"private_key":     opts.SSHKeyPairID,
	}
	params := map[string]any{
		"name":       opts.Name,
		"attributes": attrs,
	}

	if includeType {
		params["type"] = KeychainCredentialTypeSSHCredential
	}

	return params
}

// sshCredentialFromResponse converts a wire-format SSHCredentialResponse to a user-facing SSHCredential.
func sshCredentialFromResponse(resp SSHCredentialResponse) SSHCredential {
	return SSHCredential{
		ID:             resp.ID,
		Name:           resp.Name,
		Host:           resp.Host,
		Port:           resp.Port,
		Username:       resp.Username,
		RemoteHostKey:  resp.RemoteHostKey,
		ConnectTimeout: resp.ConnectTimeout,
		SSHKeyPairID:   resp.PrivateKey,
	}
}
