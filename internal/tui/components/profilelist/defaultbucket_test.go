package profilelist

import (
	"os"
	"path/filepath"
	"testing"

	appcfg "github.com/LinPr/lazys3/internal/config"
)

func writeTempConfig(t *testing.T, content string) appcfg.AWSFiles {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return appcfg.AWSFiles{ConfigFile: path}
}

func TestDefaultProfileBucket(t *testing.T) {
	// s3_bucket present: both bucket and endpoint_url are returned.
	files := writeTempConfig(t, `[default]
region = us-east-1
s3_bucket = my-bucket
endpoint_url = https://minio.example.com
`)
	bucket, endpoint := DefaultProfileBucket(files)
	if bucket != "my-bucket" {
		t.Errorf("bucket = %q, want %q", bucket, "my-bucket")
	}
	if endpoint != "https://minio.example.com" {
		t.Errorf("endpoint = %q, want %q", endpoint, "https://minio.example.com")
	}

	// No s3_bucket key: empty bucket, but endpoint_url is still surfaced.
	files = writeTempConfig(t, `[default]
region = us-east-1
endpoint_url = https://oss.aliyuncs.com
`)
	bucket, endpoint = DefaultProfileBucket(files)
	if bucket != "" {
		t.Errorf("bucket = %q, want empty", bucket)
	}
	if endpoint != "https://oss.aliyuncs.com" {
		t.Errorf("endpoint = %q, want %q", endpoint, "https://oss.aliyuncs.com")
	}

	// No [default] section: both empty.
	files = writeTempConfig(t, `[other]
s3_bucket = my-bucket
`)
	bucket, endpoint = DefaultProfileBucket(files)
	if bucket != "" || endpoint != "" {
		t.Errorf("bucket = %q, endpoint = %q, want both empty", bucket, endpoint)
	}

	// Missing config file: treated as empty, not an error.
	bucket, endpoint = DefaultProfileBucket(appcfg.AWSFiles{
		ConfigFile: filepath.Join(t.TempDir(), "nonexistent"),
	})
	if bucket != "" || endpoint != "" {
		t.Errorf("missing file: bucket = %q, endpoint = %q, want both empty", bucket, endpoint)
	}

	// Unresolved config path: both empty.
	bucket, endpoint = DefaultProfileBucket(appcfg.AWSFiles{})
	if bucket != "" || endpoint != "" {
		t.Errorf("no path: bucket = %q, endpoint = %q, want both empty", bucket, endpoint)
	}
}
