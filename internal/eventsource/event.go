package eventsource

// Event is a normalized GitHub shipping event ready for the draft generator.
// Phase 3 will add parsing functions to this package that produce Event values
// from raw GitHub JSON payloads.
type Event struct {
	Type  string // "release"
	Repo  string // "owner/repo"
	URL   string // canonical link to the release
	Title string // release name, e.g. "v0.3.0 - Notification batching"
	Body  string // release notes or PR description
	Tag   string // tag name, e.g. "v0.3.0"
}
