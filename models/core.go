package models

type ServerParameters struct {
	Port     int
	CertFile string
	KeyFile  string
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Requests struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty`
}

type Resources struct {
	Requests Requests `json:"requestsomitempty"`
}
