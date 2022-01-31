package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "mutating-webhook/models"
	server "mutating-webhook/server"
	util "mutating-webhook/utils"

	"k8s.io/client-go/kubernetes"
)

var (
	clientSet  *kubernetes.Clientset
	parameters models.ServerParameters
)

func main() {
	fmt.Printf("Starting")
	flag.IntVar(&parameters.Port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	cs, err := kubernetes.NewForConfig(util.GetKubeConfig())
	if err != nil {
		panic(err.Error())
	}
	clientSet = cs

	http.HandleFunc("/", server.HandleRoot)
	http.HandleFunc("/addresourcerequests", server.HandleMutate)

	//http.ListenAndServe("localhost:"+strconv.Itoa(parameters.Port), nil)

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(parameters.Port), parameters.CertFile, parameters.KeyFile, nil))
}
