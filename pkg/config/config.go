package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Resource struct {
	Pod                   bool `json:"pod"`
	Deployment            bool `json:"depoyment"`
	ConfigMap             bool `json:"configMap"`
	Ingress               bool `json:"ingress"`
	Job                   bool `json:"job"`
	Namespace             bool `json:"namespace"`
	PersistentVolume      bool `json:"persistentVolume"`
	ReplicaSet            bool `json:"replicaSet"`
	ReplicationController bool `json:"replicationController"`
	Secret                bool `json:"secret"`
	Service               bool `json:"service"`
}

type Config struct {
	Namespace string   `json:"namespace"`
	Resource  Resource `json:"resource"`
}

func LoadConfig() *Config {
	config := &Config{
		Namespace: "",
		Resource: Resource{
			Pod: true,
		},
	}

	file, err := os.Open(".kube-event-logger.yaml")
	if err != nil {
		fmt.Println("No config found, continuing with default")
		return config
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read config, continuing with default")
		return config
	}

	if len(b) != 0 {
		yaml.Unmarshal(b, config)
	}

	return config
}
