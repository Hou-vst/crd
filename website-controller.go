package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/Hou-vst/crd/pkg/v1"
)

var templateDir string = "./template/"
var destIpAndPort string = "localhost:8001"

func main() {
	log.Println("website-controller started.templateDir:",templateDir,",destIpAndPort:",destIpAndPort)
	for {
		resp, err := http.Get(fmt.Sprintf("http://%s/apis/extensions.example.com/v1/websites?watch=true",destIpAndPort))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var event v1.WebsiteWatchEvent
			if err := decoder.Decode(&event); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			log.Printf("Received watch event: %s: %v\n", event.Type, event.Object)

			if event.Type == "ADDED" {
				createWebsite(event.Object)
			} else if event.Type == "DELETED" {
				deleteWebsite(event.Object)
			}
		}
	}

}

func createWebsite(website v1.Website) {
	createResource(website, "api/v1", "services" ,fmt.Sprintf("%s%s",templateDir,"service-template.json"))
	createResource(website, "apis/apps/v1", "deployments", fmt.Sprintf("%s%s",templateDir,"deployment-template.json"))
	createResource(website, "apis/networking.k8s.io/v1", "ingresses", fmt.Sprintf("%s%s",templateDir,"ingress-template.json"))
}

func deleteWebsite(website v1.Website) {
	deleteResource(website, "api/v1", "services", getName(website))
	deleteResource(website, "apis/apps/v1", "deployments", getName(website))
	deleteResource(website, "apis/networking.k8s.io/v1", "ingresses", getName(website))
}

func createResource(webserver v1.Website, apiGroup string, kind string, filename string) {
	log.Printf("Creating %s with name %s in namespace %s", kind, getName(webserver), webserver.Metadata.Namespace)
	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	template := strings.Replace(string(templateBytes), "[NAME]", getName(webserver), -1)
	template = strings.Replace(template, "[NUM]", fmt.Sprint(webserver.Spec.InsNum), -1)
	template = strings.Replace(template, "[IMAGE-NAME]", webserver.Spec.Image, -1)
	template = strings.Replace(template, "[HOST-NAME]", webserver.Spec.Host, -1)

	log.Println("createResource template :", template)

	url := fmt.Sprintf("http://%s/%s/namespaces/%s/%s/",destIpAndPort ,apiGroup, webserver.Metadata.Namespace, kind)
	contentType := "application/json"
	body := strings.NewReader(template)
	resp, err := http.Post(url,contentType,body)
	log.Println("createResource url :", url)
	log.Println("createResource contentType :", contentType)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
}

func deleteResource(webserver v1.Website, apiGroup string, kind string, name string) {
	log.Printf("Deleting %s with name %s in namespace %s", kind, name, webserver.Metadata.Namespace)
	url := fmt.Sprintf("http://%s/%s/namespaces/%s/%s/%s",destIpAndPort ,apiGroup, webserver.Metadata.Namespace, kind, name)
	log.Println("deleteResource url :", url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("response Status:", resp.Status)

}

func getName(website v1.Website) string {
	return website.Metadata.Name + "-website"
}
