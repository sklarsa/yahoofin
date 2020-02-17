package yahoofin

type ServerErrorRoot struct {
	Chart ServerError `json:"chart"`
}

// ServerError is used to deserialize non-200 responses from yahoo's server
type ServerError struct {
	Result interface{}      `json:"result"`
	Error  ErrorDescription `json:"error"`
}

// ErrorDescription is used to deserialize additional error information from yahoo's server
type ErrorDescription struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
