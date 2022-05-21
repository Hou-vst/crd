package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvr = schema.GroupVersionResource{
	Group:    "extensions.example.com",
	Version:  "v1",
	Resource: "websites",
}


type Metadata struct {
	Name string 		`json:"name"`
	Namespace string	`json:"namespace"`
}

type WebsiteSpec struct {
	Host string			`json:"host"`
	Image string		`json:"image"`
	InsNum int			`json:"insNum"`
}

type Website struct {
	Metadata Metadata		`json:"metadata"`
	Spec WebsiteSpec		`json:"spec"`
}

type WebsiteWatchEvent struct {
	Type string
	Object Website
}