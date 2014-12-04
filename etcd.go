package main

import (
	"log"
	"time"
	"net/url"

	"github.com/coreos/go-etcd/etcd"
)

type EtcdStore struct {
	client    *etcd.Client
	waitIndex uint64
}

func NewEtcdStore(uri *url.URL) ConfigStore {
	urls := make([]string, 0)
	if uri.Host != "" {
		urls = append(urls, "http://"+uri.Host)
	}
	return &EtcdStore{client: etcd.NewClient(urls)}
}

func (s *EtcdStore) List(path string) (list []string) {
	resp, err := s.client.Get(path, false, true)
	if err != nil {
		log.Println("etcd:", err)
		return
	}
	if resp.Node == nil {
		return
	}
	if len(resp.Node.Nodes) == 0 {
		list = append(list, string(resp.Node.Value))
	} else {
		for _, node := range resp.Node.Nodes {
			list = append(list, string(node.Value))
		}
	}
	return
}

func (s *EtcdStore) Get(path string) string {
	resp, err := s.client.Get(path, false, false)
	if err != nil {
		log.Println("etcd:", err)
		return ""
	}
	if resp.Node == nil {
		return ""
	}
	return string(resp.Node.Value)
}

func (s *EtcdStore) Watch(path string) {
	resp, err := s.client.Watch(path, s.waitIndex, true, nil, nil)
	if err == nil {
		s.waitIndex = resp.EtcdIndex + 1
	} else {
		etcdError, ok := err.(*etcd.EtcdError)
		if !ok {
			log.Println("etcd:", err)
			// Let's not slam the etcd server in the event that we know
			// an unexpected error occurred.
			time.Sleep(time.Second)
			return
		}

		switch etcdError.ErrorCode {
		case 401: //etcdError.EcodeEventIndexCleared:
			// This is racy, but adding one to the last known index
			// will help get this watcher back into the range of
			// etcd's internal event history
			s.waitIndex = s.waitIndex + 1
		default:
			log.Println("etcd:", err)
			time.Sleep(time.Second)
		}
	}
}

