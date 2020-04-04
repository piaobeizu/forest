package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"strings"
	"time"

	"github.com/piaobeizu/forest"
	"github.com/prometheus/common/log"
)

const (
	DefaultEndpoints   = "127.0.0.1:2379"
	DefaultHttpAddress = ":2856"
	DefaultDialTimeout = 5
	DefaultDbUrl       = "root:GypassBDN^1024@tcp(127.0.0.1:13306)/rootcloud_trident?charset=utf8"
	DefaultEtcdCert    = "/etc/kubernetes/pki/etcd/ca.crt"
	DefaultEtcdKey     = "/etc/kubernetes/pki/etcd/ca.key"
)

func main() {

	ip := forest.GetLocalIpAddress()
	if ip == "" {
		log.Fatal("has no get the ip address")

	}

	endpoints := flag.String("etcd-endpoints", DefaultEndpoints, "etcd endpoints")
	httpAddress := flag.String("http-address", DefaultHttpAddress, "http address")
	etcdCertFile := flag.String("etcd-cert", DefaultEtcdCert, "etcd-cert file")
	etcdKeyFile := flag.String("etcd-key", DefaultEtcdKey, "etcd-key file")
	etcdDialTime := flag.Int64("etcd-dailtimeout", DefaultDialTimeout, "etcd dailtimeout")
	help := flag.String("help", "", "forest help")
	dbUrl := flag.String("db-url", DefaultDbUrl, "db-url for mysql")
	flag.Parse()
	if *help != "" {
		flag.Usage()
		return
	}

	endpoint := strings.Split(*endpoints, ",")
	dialTime := time.Duration(*etcdDialTime) * time.Second

	var (
		etcd *forest.Etcd
		err  error
		cfg  *tls.Config
	)
	if etcdCertFile != nil && etcdKeyFile != nil {
		cert, err := ioutil.ReadFile(*etcdCertFile)
		if err != nil {
			log.Fatal(err)
		}

		key, err := ioutil.ReadFile(*etcdKeyFile)
		if err != nil {
			log.Fatal(err)
		}

		certificate, err := tls.X509KeyPair(cert, key)
		if err != nil {
			log.Fatal(err)
		}
		cfg = &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{certificate},
		}
	}
	etcd, err = forest.NewEtcd(endpoint, dialTime, cfg)
	if err != nil {
		log.Fatal(err)
	}

	node, err := forest.NewJobNode(ip, etcd, *httpAddress, *dbUrl)
	if err != nil {

		log.Fatal(err)
	}

	node.Bootstrap()
}
