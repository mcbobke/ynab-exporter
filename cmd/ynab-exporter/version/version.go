// Package version provides buildtime-populated vars for meta-metrics.
package version

var (
	// BuildVersion is populated at compile time with the tagged version of the exporter.
	BuildVersion string

	// BuildTime is populated at compile time with the timestamp of the time this version of the exporter was built.
	BuildTime string
)
