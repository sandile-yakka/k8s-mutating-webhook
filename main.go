package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"flag"

	"k8s.io/api/admission/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ServerParameters struct {
	port     int
	certFile string
	keyFile  string
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

var (
	clientSet             *kubernetes.Clientset
	parameters            ServerParameters
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func main() {

	flag.IntVar(&parameters.port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	cs, err := kubernetes.NewForConfig(getKubeConfig())
	if err != nil {
		panic(err.Error())
	}
	clientSet = cs
	// test()

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/AddResourceRequests", HandleMutate)

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(parameters.port), parameters.certFile, parameters.keyFile, nil))

}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HandleRoot!"))
}

func HandleMutate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("could not deserialize request: %v", err)
	} else if admissionReviewReq.Request == nil {
		w.WriteHeader(http.StatusBadRequest)
		errors.New("malformed admission review: request is nil")
	}

	fmt.Printf("Type: %v \t Event: %v \t Name: %v \n", admissionReviewReq.Request.Kind, admissionReviewReq.Request.Operation, admissionReviewReq.Request.Name)

	var pod apiv1.Pod
	err = json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod)
	if err != nil {
		fmt.Errorf("could not unmarshal pod on admission request: %v", err)
	}

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
		},
	}

	patchBytes, err := json.Marshal(patchResourceRequests(pod))

	if err != nil {
		fmt.Errorf("could not marshal JSON patch: %v", err)
	}

	admissionReviewResponse.Response.Patch = patchBytes

	bytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		fmt.Errorf("marshalling repsonse: %v", err)
	}

	w.Write(bytes)
}

func patchResourceRequests(pod apiv1.Pod) []patchOperation {
	var patches []patchOperation

	Requests := pod.Spec.Containers[0].Resources.Requests
	cpu := resource.NewScaledQuantity(250, resource.Mega)
	mega := resource.NewScaledQuantity(100, resource.Mega)

	Requests.Cpu().Add(*cpu)
	Requests.Memory().Add(*mega)

	patches = append(patches, patchOperation{
		Op:    "add",
		Path:  "/spec/containers/0/resources/requests",
		Value: Requests,
	})

	return patches
}

func getKubeConfig() *rest.Config {
	var config *rest.Config
	kubeConfigFilePath := os.Getenv("KUBECONFIG")
	useKubeConfig, err := strconv.ParseBool(os.Getenv("USE_KUBECONFIG"))

	if err != nil {
		panic(err.Error())
	}

	if !useKubeConfig {
		c, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		config = c
	} else {
		kubeconfig := kubeConfigFilePath

		if kubeConfigFilePath == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}
		}

		c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		config = c
	}
	return config
}
