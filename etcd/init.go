package etcd

var Store *store

func Initialize(){
	Store = NewStore(60)
}
