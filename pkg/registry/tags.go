package registry

import (
	"context"
	"fmt"
	"strings"

	"oras.land/oras-go/v2/registry/remote"
)

func ParseOCIURL(ociURL string) (host, repo string, err error) {
	if !strings.HasPrefix(ociURL, "oci://") {
		return "", "", fmt.Errorf("URL must start with oci://, got %q", ociURL)
	}

	rest := strings.TrimPrefix(ociURL, "oci://")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) < 2 || parts[1] == "" {
		return "", "", fmt.Errorf("OCI URL must include a repository path: %q", ociURL)
	}

	host = parts[0]
	repo = parts[1]

	if host == "docker.io" {
		host = "registry-1.docker.io"
	}

	return host, repo, nil
}

func ListTags(ociURL string, plainHTTP bool) ([]string, error) {
	host, repoPath, err := ParseOCIURL(ociURL)
	if err != nil {
		return nil, err
	}

	ref := host + "/" + repoPath
	repo, err := remote.NewRepository(ref)
	if err != nil {
		return nil, fmt.Errorf("connecting to registry: %w", err)
	}
	repo.PlainHTTP = plainHTTP

	var tags []string
	err = repo.Tags(context.Background(), "", func(t []string) error {
		tags = append(tags, t...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}

	return tags, nil
}
