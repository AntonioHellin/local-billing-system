package frontend

import "embed"

//go:embed dist
// FS contains the built React application static files.
var FS embed.FS
