package db

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

type Etcd struct {
	Ctx context.Context
	KV  clientv3.KV
}

func NewEtcdClient(address string) *Etcd {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{address},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	return &Etcd{Ctx: ctx, KV: kv}
}
