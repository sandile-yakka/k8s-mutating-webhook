package utils

import (
	"os"
	"path/filepath"
	"strconv"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetKubeConfig() *rest.Config {
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
