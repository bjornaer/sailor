package db

// DBClient interface just handles strings *but* I would totally use generics here
type DBClient interface {
	Put(k, v string) error
	Get(k string) (string, error)
}

// InitDBClient takes variadic args, if none is given it uses an in memory Map
// if "etcd" is passed, a second argument for the address is required
func InitDBClient(clientType ...string) DBClient {
	return NewMapDBClient()
	// if len(clientType) == 0 {
	// 	// default use in memory map
	// 	return NewMapDBClient()
	// }
	// // if not default use whatever came from variable first
	// dbClientType := clientType[0]
	// if dbClientType == "etcd" {
	// 	address := clientType[1]
	// 	return NewEtcdClient(address)
	// } else {
	// 	// for anything else we go back to the in memory map
	// 	// Ideally we mighjjt have more DB type of options and handle those
	// 	return NewMapDBClient()
	// }
}
