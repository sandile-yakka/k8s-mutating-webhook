package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	models "mutating-webhook/models"

	"k8s.io/api/admission/v1beta1"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

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

	var deployment apps.Deployment
	err = json.Unmarshal(admissionReviewReq.Request.Object.Raw, &deployment)
	if err != nil {
		fmt.Errorf("could not unmarshal pod on admission request: %v", err)
	}

	patchType := "JSONPatch"

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:       admissionReviewReq.Request.UID,
			Allowed:   true,
			PatchType: (*v1beta1.PatchType)(&patchType),
		},
	}

	patchBytes, err := json.Marshal(PatchResourceRequests())
	if err != nil {
		fmt.Errorf("could not marshal JSON patch: %v", err)
	}

	admissionReviewResponse.Response.Patch = patchBytes
	bytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		fmt.Errorf("marshalling repsonse: %v", err)
	}

	fmt.Print(string(bytes))

	w.Write(bytes)
}

func PatchResourceRequests() []models.PatchOperation {
	var patches []models.PatchOperation
	resources := models.Resources{
		Requests: models.Requests{
			Cpu:    "200m",
			Memory: "100Mi",
		},
	}

	patches = append(patches, models.PatchOperation{
		Op:    "add",
		Path:  "/spec/template/spec/containers/0/resources",
		Value: resources,
	})
	return patches
}
