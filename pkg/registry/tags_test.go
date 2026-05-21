package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type tagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func mockRegistry(t *testing.T, repo string, tags []string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/"+repo+"/tags/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tagsResponse{Name: repo, Tags: tags})
	})
	return httptest.NewServer(mux)
}

func TestListTags(t *testing.T) {
	wantTags := []string{"v1.0.0", "v1.1.0", "v2.0.0"}
	srv := mockRegistry(t, "envoyproxy/gateway-helm", wantTags)
	defer srv.Close()

	ociURL := "oci://" + srv.Listener.Addr().String() + "/envoyproxy/gateway-helm"

	got, err := ListTags(ociURL, true)
	if err != nil {
		t.Fatalf("ListTags: %v", err)
	}

	if len(got) != len(wantTags) {
		t.Fatalf("got %d tags, want %d", len(got), len(wantTags))
	}

	for i, tag := range got {
		if tag != wantTags[i] {
			t.Errorf("tag[%d] = %q, want %q", i, tag, wantTags[i])
		}
	}
}

func TestListTagsEmpty(t *testing.T) {
	srv := mockRegistry(t, "empty/chart", []string{})
	defer srv.Close()

	ociURL := "oci://" + srv.Listener.Addr().String() + "/empty/chart"

	got, err := ListTags(ociURL, true)
	if err != nil {
		t.Fatalf("ListTags: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("got %d tags, want 0", len(got))
	}
}

func TestParseOCIURL(t *testing.T) {
	tests := []struct {
		input    string
		wantHost string
		wantRepo string
		wantErr  bool
	}{
		{
			input:    "oci://docker.io/envoyproxy/gateway-helm",
			wantHost: "registry-1.docker.io",
			wantRepo: "envoyproxy/gateway-helm",
		},
		{
			input:    "oci://ghcr.io/some-org/some-chart",
			wantHost: "ghcr.io",
			wantRepo: "some-org/some-chart",
		},
		{
			input:    "oci://quay.io/jetstack/cert-manager",
			wantHost: "quay.io",
			wantRepo: "jetstack/cert-manager",
		},
		{
			input:    "oci://myregistry.example.com:5000/charts/app",
			wantHost: "myregistry.example.com:5000",
			wantRepo: "charts/app",
		},
		{
			input:   "https://not-oci.example.com/chart",
			wantErr: true,
		},
		{
			input:   "oci://docker.io",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			host, repo, err := ParseOCIURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseOCIURL: %v", err)
			}
			if host != tt.wantHost {
				t.Errorf("host = %q, want %q", host, tt.wantHost)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}
