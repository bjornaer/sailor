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

func (e *Etcd) Get(k string) (string, error) {
	v, err := e.KV.Get(e.Ctx, k)
	if err != nil {
		return "", err
	}
	return string(v.Kvs[0].Value), nil
}

func (e *Etcd) Put(k, v string) error {
	_, err := e.KV.Put(e.Ctx, k, v)
	if err != nil {
		return err
	}
	return nil
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
