package response

const (
	// CodeInternal indicates the server failed to complete a request safely.
	CodeInternal = "internal_error"
	// CodeNotFound indicates the requested resource does not exist.
	CodeNotFound = "not_found"
	// CodeServiceUnavailable indicates a dependency is temporarily unavailable.
	CodeServiceUnavailable = "service_unavailable"
	// CodeMethodNotAllowed indicates the route exists but does not allow the method.
	CodeMethodNotAllowed = "method_not_allowed"
)
