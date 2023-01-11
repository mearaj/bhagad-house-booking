package frontend

const ApiURLKey = "API_URL"
const StaticFolderKey = "STATIC_FOLDER"
const InnerPortKey = "INNER_PORT"
const DefaultApiURL = "http://localhost:8001"
const DefaultStaticFolder = "dist"
const DefaultInnerPort = "8080"

type Config struct {
	ApiURL       string
	StaticFolder string
	InnerPort    string
}
