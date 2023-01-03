package frontend

import (
	"syscall/js"
)

// LoadConfig looks for env API_URL. (if empty, then it defaults to DefaultApiURL)
func LoadConfig() (config Config) {
	apiURL := js.Global().Get(ApiURLKey)
	if apiURL.IsUndefined() || apiURL.IsNull() {
		config.ApiURL = DefaultApiURL
		return config
	}
	config.ApiURL = apiURL.String()
	return config
}
