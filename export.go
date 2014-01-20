package gorest

import (
	"fmt"
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
	Endpoints    map[string]EndPointDescription
}

type EndPointDescription struct {
	Signature    string
	RawSignature string
	Method       string
	FunctionName string
	Returns      string
	PathParams   []EndPointParam
	QueryParams  []EndPointParam
}

type EndPointParam struct {
	Name string
	Type string
	Kind string
}

func GetServiceDescription(root_path string) *ServiceDescription {

	// Create a valid descriptor
	s := newServiceDescription()

	// The root path will almost certainly contain
	// slashes, which aren't that useful to some
	// of the internal storage types.
	trimmed_root_path := strings.Trim(root_path, "/")

	// Set root path.
	// Apply a leading slash to the root,
	// and remove any trailing slashes
	s.RootPath = "/" + trimmed_root_path

	// Endpoints are stored as a single array, and
	// need to be matched to services. The first step
	// is to collate them into a map of URL names.
	endpoints := make(map[string]endPointStruct)

	for _, value := range _manager().endpoints {
		key := strings.Replace(value.signiture, trimmed_root_path, "", 1)
		endpoints[key] = value
	}

	for _, v := range _manager().serviceTypes {

		// Remove root path from relative path
		path := "/" + strings.Trim(strings.Replace(v.root, s.RootPath, "", 1), "/")

		// Assign to the array
		s.Services[path] = serviceMetaDataToPublicService(v)

		// Find any associated endpoints
		for key, ep := range endpoints {
			if index := strings.Index(key, path); index == 0 {

				signature := strings.Replace(key, path, "", 1)

				// Create a new public endpoind description
				s.Services[path].Endpoints[signature] = endpointMetaDataToPublicEndpoind(ep, signature)
			}
		}
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

	out.Endpoints = make(map[string]EndPointDescription)

	return
}

func endpointMetaDataToPublicEndpoind(in endPointStruct, sig string) EndPointDescription {

	e := EndPointDescription{}

	e.RawSignature = sig
	e.Signature = strings.Split(sig, "?")[0]
	e.Method = in.requestMethod
	e.Returns = in.outputType

	// Generate a unique method name
	{
		packages := strings.Split(in.parentTypeName, "/")
		type_name := packages[len(packages)-1]
		e.FunctionName = fmt.Sprintf("%s.%s", type_name, in.methodName)
	}

	// Path params
	if len(in.params) > 0 {
		for i := 0; i < len(in.params); i++ {
			e.PathParams = append(e.PathParams, makePublicEnpointParamSet(&e, in.params[i], "path"))
		}
	}

	if len(in.queryParams) > 0 {
		for i := 0; i < len(in.queryParams); i++ {
			e.PathParams = append(e.PathParams, makePublicEnpointParamSet(&e, in.queryParams[i], "query"))
		}
	}

	return e
}

func makePublicEnpointParamSet(e *EndPointDescription, in param, kind string) (out EndPointParam) {

	// Replace the {name:type} format with {name}
	e.Signature = strings.Replace(
		e.Signature,
		fmt.Sprintf("{%s:%s}", in.name, in.typeName),
		fmt.Sprintf("{%s}", strings.ToLower(in.name)),
		1)

	// Store the details
	out.Name = strings.ToLower(in.name)
	out.Type = in.typeName
	out.Kind = kind

	return
}
