package cmd

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

type mockRunner struct {
	lastArgs []string
}

func (m *mockRunner) Run(args []string, stdout, stderr io.Writer) error {
	m.lastArgs = args
	return nil
}

func setupDelegationTest(t *testing.T) (*mockRunner, func(args ...string) error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "oci-bookmarks.yaml")
	setBookmarkPath(path)

	mock := &mockRunner{}
	setRunner(mock)
	t.Cleanup(func() { setRunner(execRunner{}) })

	root := New()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)

	run := func(args ...string) error {
		root.SetArgs(args)
		return root.Execute()
	}

	_ = run("add", "envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")

	return mock, run
}

func assertArgs(t *testing.T, got []string, wantParts ...string) {
	t.Helper()
	joined := strings.Join(got, " ")
	for _, part := range wantParts {
		if !strings.Contains(joined, part) {
			t.Errorf("args %v missing %q", got, part)
		}
	}
}

func TestValuesCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("values", "envoy-gw", "--version", "1.3.0"); err != nil {
		t.Fatalf("values: %v", err)
	}

	assertArgs(t, mock.lastArgs, "show", "values", "oci://docker.io/envoyproxy/gateway-helm", "--version", "1.3.0")
}

func TestShowCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("show", "envoy-gw"); err != nil {
		t.Fatalf("show: %v", err)
	}

	assertArgs(t, mock.lastArgs, "show", "chart", "oci://docker.io/envoyproxy/gateway-helm")
}

func TestInstallCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("install", "envoy-gw", "my-release", "--version", "1.3.0", "--namespace", "envoy"); err != nil {
		t.Fatalf("install: %v", err)
	}

	assertArgs(t, mock.lastArgs, "install", "my-release", "oci://docker.io/envoyproxy/gateway-helm", "--version", "1.3.0", "--namespace", "envoy")
}

func TestUpgradeCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("upgrade", "envoy-gw", "my-release", "--version", "1.4.0"); err != nil {
		t.Fatalf("upgrade: %v", err)
	}

	assertArgs(t, mock.lastArgs, "upgrade", "my-release", "oci://docker.io/envoyproxy/gateway-helm", "--version", "1.4.0")
}

func TestPullCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("pull", "envoy-gw", "--version", "1.3.0"); err != nil {
		t.Fatalf("pull: %v", err)
	}

	assertArgs(t, mock.lastArgs, "pull", "oci://docker.io/envoyproxy/gateway-helm", "--version", "1.3.0")
}

func TestTemplateCommand(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("template", "envoy-gw", "my-release", "--version", "1.3.0"); err != nil {
		t.Fatalf("template: %v", err)
	}

	assertArgs(t, mock.lastArgs, "template", "my-release", "oci://docker.io/envoyproxy/gateway-helm", "--version", "1.3.0")
}

func TestDelegationBookmarkNotFound(t *testing.T) {
	_, run := setupDelegationTest(t)

	if err := run("values", "nonexistent"); err == nil {
		t.Fatal("expected error for unknown bookmark")
	}
}

func TestInstallPassthroughFlags(t *testing.T) {
	mock, run := setupDelegationTest(t)

	if err := run("install", "envoy-gw", "my-release", "--version", "1.3.0", "--set", "foo=bar", "--values", "custom.yaml"); err != nil {
		t.Fatalf("install with passthrough: %v", err)
	}

	assertArgs(t, mock.lastArgs, "--set", "foo=bar", "--values", "custom.yaml")
}
