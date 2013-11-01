package gorest

import (
	"strings"
)

type ServiceDescription struct {
	RootPath string
	Services map[string]Service
}

type Service struct {
	Description  string
	ConsumesMime string
	ProducesMime string
}

func GetServiceDescription(root_path string) *ServiceDescription {

	// Create a valid descriptor
	s := newServiceDescription()

	// Set root path.
	// Apply a leading slash to the root,
	// and remove any trailing slashes
	s.RootPath = "/" + strings.Trim(root_path, "/")

	for _, v := range _manager().serviceTypes {

		// Remove root path from relative path
		path := "/" + strings.Trim(strings.Replace(v.root, s.RootPath, "", 1), "/")

		// Assign to the array
		s.Services[path] = serviceMetaDataToPublicService(v)
	}

	return s
}

func newServiceDescription() *ServiceDescription {
	s := new(ServiceDescription)

	s.Services = make(map[string]Service)

	return s
}

func serviceMetaDataToPublicService(in serviceMetaData) (out Service) {

	out.ConsumesMime = in.consumesMime
	out.ProducesMime = in.producesMime
	out.Description = in.description

	return
}
