package transport

type RemoteService interface {
	GetServiceName() string
	GetServiceUrl() []string
}

type HttpRemoteService struct {
	ServiceName string
	Url         []string
}

func NewHttpRemoteService(serviceName string, urls []string) RemoteService {
	return &HttpRemoteService{
		ServiceName: serviceName,
		Url:         urls,
	}
}

func (r *HttpRemoteService) GetServiceName() string {
	return r.ServiceName
}

func (r *HttpRemoteService) GetServiceUrl() []string {
	return r.Url
}
