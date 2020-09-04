package etcd

// Creating a cluster connection store for global access
var ClusterConnection *ConnectionStore

// Initializing the cluster connection store with 15 minutes of ttl period
func Initialize(){
	ClusterConnection = New(15)
}
