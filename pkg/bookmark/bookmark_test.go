package bookmark

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "oci-bookmarks.yaml")
}

func TestAddAndGet(t *testing.T) {
	s := NewStore(tempStorePath(t))

	if err := s.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	b, err := s.Get("envoy-gw")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if b.Name != "envoy-gw" {
		t.Errorf("Name = %q, want %q", b.Name, "envoy-gw")
	}
	if b.URL != "oci://docker.io/envoyproxy/gateway-helm" {
		t.Errorf("URL = %q, want %q", b.URL, "oci://docker.io/envoyproxy/gateway-helm")
	}
}

func TestAddDuplicate(t *testing.T) {
	s := NewStore(tempStorePath(t))

	if err := s.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm"); err != nil {
		t.Fatalf("first Add: %v", err)
	}
	err := s.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")
	if err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestAddInvalidURL(t *testing.T) {
	s := NewStore(tempStorePath(t))

	err := s.Add("bad", "https://not-oci.example.com/chart")
	if err == nil {
		t.Fatal("expected error for non-oci URL, got nil")
	}
}

func TestRemove(t *testing.T) {
	s := NewStore(tempStorePath(t))

	_ = s.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")

	if err := s.Remove("envoy-gw"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	_, err := s.Get("envoy-gw")
	if err == nil {
		t.Fatal("expected error after removal, got nil")
	}
}

func TestRemoveNotFound(t *testing.T) {
	s := NewStore(tempStorePath(t))

	err := s.Remove("nonexistent")
	if err == nil {
		t.Fatal("expected error removing nonexistent bookmark, got nil")
	}
}

func TestList(t *testing.T) {
	s := NewStore(tempStorePath(t))

	_ = s.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")
	_ = s.Add("cert-manager", "oci://quay.io/jetstack/cert-manager")

	items, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("List len = %d, want 2", len(items))
	}
}

func TestListEmpty(t *testing.T) {
	s := NewStore(tempStorePath(t))

	items, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("List len = %d, want 0", len(items))
	}
}

func TestPersistence(t *testing.T) {
	path := tempStorePath(t)

	s1 := NewStore(path)
	_ = s1.Add("envoy-gw", "oci://docker.io/envoyproxy/gateway-helm")

	s2 := NewStore(path)
	b, err := s2.Get("envoy-gw")
	if err != nil {
		t.Fatalf("Get from second store instance: %v", err)
	}
	if b.URL != "oci://docker.io/envoyproxy/gateway-helm" {
		t.Errorf("URL = %q after reload", b.URL)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewStore(tempStorePath(t))

	_, err := s.Get("nope")
	if err == nil {
		t.Fatal("expected error for missing bookmark, got nil")
	}
}

func TestFileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "subdir", "bookmarks.yaml")

	s := NewStore(path)
	if err := s.Add("test", "oci://example.com/chart"); err != nil {
		t.Fatalf("Add with missing parent dir: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
