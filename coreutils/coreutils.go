package coreutils
	
type Request struct {
	Method string
	Path string
	Version string
	Headers map[string]string
	Body []byte
}

type Respone struct {
	Status_code int
	Headers map[string]string
	Body []byte
}
