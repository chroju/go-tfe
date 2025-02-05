package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ StateVersionOutputs = (*stateVersionOutputs)(nil)

// State version outputs are the output values from a Terraform state file.
// They include the name and value of the output, as well as a sensitive boolean
// if the value should be hidden by default in UIs.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/state-version-outputs.html
type StateVersionOutputs interface {
	Read(ctx context.Context, outputID string) (*StateVersionOutput, error)
	ReadCurrent(ctx context.Context, workspaceID string) (*StateVersionOutputsList, error)
}

// stateVersionOutputs implements StateVersionOutputs.
type stateVersionOutputs struct {
	client *Client
}

// StateVersionOutput represents a State Version Outputs
type StateVersionOutput struct {
	ID        string      `jsonapi:"primary,state-version-outputs"`
	Name      string      `jsonapi:"attr,name"`
	Sensitive bool        `jsonapi:"attr,sensitive"`
	Type      string      `jsonapi:"attr,type"`
	Value     interface{} `jsonapi:"attr,value"`
}

// ReadCurrent reads the current state version outputs for the specified workspace
func (s *stateVersionOutputs) ReadCurrent(ctx context.Context, workspaceID string) (*StateVersionOutputsList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/current-state-version-outputs", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	so := &StateVersionOutputsList{}
	err = s.client.do(ctx, req, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}

// Read a State Version Output
func (s *stateVersionOutputs) Read(ctx context.Context, outputID string) (*StateVersionOutput, error) {
	if !validStringID(&outputID) {
		return nil, ErrInvalidOutputID
	}

	u := fmt.Sprintf("state-version-outputs/%s", url.QueryEscape(outputID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	so := &StateVersionOutput{}
	err = s.client.do(ctx, req, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}
