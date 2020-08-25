package model

type Cluster struct {

	// ID is a unique identifier for the cluster in the watcher application
	ID string `json:"id"`

	// Name is a user-given name for the cluster
	Name string `json:"name" validate:"required"`

	// Endpoints is a list of URLs of nodes
	Endpoints []string `json:"endpoints"`

	// Username is a user name for authentication.
	Username string `json:"username"`

	// Password is a password for authentication.
	Password string `json:"password"`

	// TLS indicates whether to use tls certificates for authentication
	TLS bool `json:"tls"`

	// ServerName ensures the cert matches the given host in case of discovery / virtual hosting
	ServerName string `json:"server_name"`

	// CertFile is Certificate file name
	CertFile string `json:"cert_file"`

	// KeyFile is key file name
	KeyFile string `json:"key_file"`

	// TrustedCAFile is Certificate Authority file name
	TrustedCAFile string `json:"trusted_ca_file"`

	// CreationTime indicates the time when the cluster profile was created
	CreationTime int64 `json:"creation_time"`
}
