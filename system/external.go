package system

// External represents any external dependencies or modules needed by a module.
type External struct {
	// Name is the identifier of the external dependency.
	Name string `prompt:"What is the name of this external dependency?" default:"External"`

	// Version specifies the version of the external dependency.
	Version string `prompt:"What version of this external dependency do you need?" default:"latest"`

	// DownloadURL specifies the URL to download the external dependency from.
	DownloadURL string
}
