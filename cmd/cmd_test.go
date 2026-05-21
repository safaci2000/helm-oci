package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestCmd(t *testing.T) (*bytes.Buffer, *bytes.Buffer, func(args ...string) error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "oci-bookmarks.yaml")

	root := New()
	root.SetArgs([]string{})

	setBookmarkPath(path)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)

	run := func(args ...string) error {
		stdout.Reset()
		stderr.Reset()
		root.SetArgs(args)
		return root.Execute()
	}

	return &stdout, &stderr, run
}

func TestAddCommand(t *testing.T) {
	_, _, run := setupTestCmd(t)

	if err := run("add", "envoy-gw", "oci://docker.io/envoyproxy/gateway-helm"); err != nil {
		t.Fatalf("add: %v", err)
	}
}

func TestAddCommandMissingArgs(t *testing.T) {
	_, _, run := setupTestCmd(t)

	if err := run("add", "envoy-gw"); err == nil {
		t.Fatal("expected error with missing URL arg")
	}
}

func TestAddCommandInvalidURL(t *testing.T) {
	_, _, run := setupTestCmd(t)

	if err := run("add", "envoy-gw", "https://not-oci.com/chart"); err == nil {
		t.Fatal("expected error for non-oci URL")
	}
}

func TestListCommand(t *testing.T) {
	stdout, _, run := setupTestCmd(t)

	_ = run("add", "envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")
	_ = run("add", "cert-mgr", "oci://quay.io/jetstack/cert-manager")

	if err := run("list"); err != nil {
		t.Fatalf("list: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "envoy-gw") {
		t.Errorf("list output missing envoy-gw:\n%s", out)
	}
	if !strings.Contains(out, "cert-mgr") {
		t.Errorf("list output missing cert-mgr:\n%s", out)
	}
	if !strings.Contains(out, "oci://docker.io/envoyproxy/gateway-helm") {
		t.Errorf("list output missing URL:\n%s", out)
	}
}

func TestListCommandEmpty(t *testing.T) {
	stdout, _, run := setupTestCmd(t)

	if err := run("list"); err != nil {
		t.Fatalf("list: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "No bookmarks") {
		t.Errorf("expected 'No bookmarks' message, got:\n%s", out)
	}
}

func TestRemoveCommand(t *testing.T) {
	stdout, _, run := setupTestCmd(t)

	_ = run("add", "envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")

	if err := run("remove", "envoy-gw"); err != nil {
		t.Fatalf("remove: %v", err)
	}

	if err := run("list"); err != nil {
		t.Fatalf("list after remove: %v", err)
	}

	if strings.Contains(stdout.String(), "envoy-gw") {
		t.Error("envoy-gw still present after removal")
	}
}

func TestRemoveCommandNotFound(t *testing.T) {
	_, _, run := setupTestCmd(t)

	if err := run("remove", "nonexistent"); err == nil {
		t.Fatal("expected error removing nonexistent bookmark")
	}
}
