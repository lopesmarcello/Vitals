package viewhelpers

// GetStatusColor returns the Tailwind classes for a given HTTP status code
func GetStatusColor(code int) string {
	if code >= 200 && code < 300 {
		return "bg-green-500/10 text-green-400 border-green-500/20"
	}
	if code >= 300 && code < 400 {
		return "bg-yellow-500/10 text-yellow-400 border-yellow-500/20"
	}
	if code >= 400 || code == 0 {
		return "bg-red-500/10 text-red-400 border-red-500/20"
	}
	return "bg-slate-500/10 text-slate-400 border-slate-500/20"
}
