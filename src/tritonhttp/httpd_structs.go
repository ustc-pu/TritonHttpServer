package tritonhttp

type HttpServer	struct {
	ServerPort	string
	DocRoot		string
	MIMEPath	string
	MIMEMap		map[string]string
}

type HttpResponseHeader struct {
	// Add any fields required for the response here
	version string
	code string
	description string
	responseHeaderMap map[string]string
	body string
}

type HttpRequestHeader struct {
	// Add any fields required for the request here
	method string
	url string
	version string
	requestHeaderMap map[string]string
}