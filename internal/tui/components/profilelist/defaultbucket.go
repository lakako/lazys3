package profilelist

import (
	"log"

	appcfg "github.com/LinPr/lazys3/internal/config"
)

const (
	// DefaultProfileName is the shared-config section lazys3 inspects for
	// the startup auto-entry key.
	DefaultProfileName = "default"

	// s3BucketKey is a lazys3 extension to the shared config format: when
	// set under [default] in ~/.aws/config, lazys3 skips the profile
	// picker and bucket list and opens that bucket directly on startup.
	s3BucketKey = "s3_bucket"

	// endpointURLKey is the standard shared-config endpoint override; it
	// is read alongside s3BucketKey so non-AWS services (OSS, OBS, MinIO,
	// ...) connect through the right endpoint in the auto-entered bucket.
	endpointURLKey = "endpoint_url"
)

// DefaultProfileBucket returns the s3_bucket configured on the [default]
// section of the resolved shared config file, plus that section's
// endpoint_url ("" when unset). A missing file, section or key yields empty
// strings so callers keep the normal profile-picker flow — the auto-entry
// is purely additive and never turns a malformed shared config into a hard
// startup error.
func DefaultProfileBucket(files appcfg.AWSFiles) (bucket, endpointURL string) {
	if files.ConfigFile == "" {
		return "", ""
	}
	// Only the shared config file is scanned: s3_bucket is a lazys3
	// extension and the credentials file has no business carrying it.
	// loadIniFiles treats an unreadable file as empty, so a missing
	// ~/.aws/config simply yields no key.
	sections, err := loadIniFiles(fileSlice(files.ConfigFile))
	if err != nil {
		log.Println("read s3_bucket from shared config failed:", err)
		return "", ""
	}
	if !sections.HasSection(DefaultProfileName) {
		return "", ""
	}
	section, _ := sections.GetSection(DefaultProfileName)
	if section.Has(s3BucketKey) {
		bucket = section.String(s3BucketKey)
	}
	if section.Has(endpointURLKey) {
		endpointURL = section.String(endpointURLKey)
	}
	return bucket, endpointURL
}
