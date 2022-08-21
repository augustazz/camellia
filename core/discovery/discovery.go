package discovery

import (
	"github.com/augustazz/camellia/config"
	"github.com/augustazz/camellia/constants"
	"github.com/augustazz/camellia/core/transport"
)

var discovery Discovery

type DiscoveryType uint8

const (
	Local DiscoveryType = iota
	Etcd
)

func GetDiscoveryType(t int) DiscoveryType {
	return DiscoveryType(t)
}

type Discovery interface {
	Startup() error
	ServiceList() []transport.RemoteService
	GetServices(serviceName string) transport.RemoteService
}

func NewDiscovery(t DiscoveryType) Discovery {
	return discovery
}

func GetDiscovery() Discovery {
	return discovery
}

type LocalConfigDiscovery struct {
}

func (d *LocalConfigDiscovery) Startup() error {
	services := config.GetServiceConfig()
	if len(services) == 0 {
		return constants.ConfigServicesNotFound
	}
	return nil
}

func (d *LocalConfigDiscovery) GetServices(serviceName string) transport.RemoteService {
	services := config.GetServiceConfig()
	for _, v := range services {
		if v.Name == serviceName {
			url := v.Url
			if len(url) > 0 {
				return transport.NewHttpRemoteService(v.Name, v.Url)
			}
		}
	}
	return nil
}

func (d *LocalConfigDiscovery) ServiceList() []transport.RemoteService {
	services := config.GetServiceConfig()
	r := make([]transport.RemoteService, 0)
	for _, v := range services {
		r = append(r, transport.NewHttpRemoteService(v.Name, v.Url))
	}
	return r
}

type EtcdServiceDiscovery struct {
}

func (d *EtcdServiceDiscovery) Startup() error {
	return nil
}

func (d *EtcdServiceDiscovery) GetServices(serviceName string) transport.RemoteService {

	return nil
}

func (d *EtcdServiceDiscovery) ServiceList() []transport.RemoteService {
	return nil
}
