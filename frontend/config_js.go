package frontend

import (
	"syscall/js"
)

// LoadConfig looks for env API_URL. (if empty, then it defaults to DefaultApiURL)
func LoadConfig() (config Config) {
	apiURL := js.Global().Get(ApiURLKey)
	if apiURL.IsUndefined() || apiURL.IsNull() {
		config.ApiURL = DefaultApiURL
	} else {
		config.ApiURL = apiURL.String()
	}
	staticFolder := js.Global().Get(StaticFolderKey)
	if staticFolder.IsUndefined() || staticFolder.IsNull() {
		config.StaticFolder = DefaultStaticFolder
	} else {
		config.StaticFolder = staticFolder.String()
	}
	innerPort := js.Global().Get(InnerPortKey)
	if innerPort.IsUndefined() || innerPort.IsNull() {
		config.InnerPort = DefaultInnerPort
	} else {
		config.InnerPort = innerPort.String()
	}
	return config
}
