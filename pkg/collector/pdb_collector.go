package collector

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type PDBInfo struct {
	Name               string `json:"name"`
	MinAvailable       string `json:"minavailable"`
	MaxUnavailable     string `json:"maxunavailable"`
	DisruptionsAllowed int32  `json:"disruptionsallowed"`
}

// PDBCollector defines a Pod disruption Budget Collector struct
type PDBCollector struct {
	kubeconfig *restclient.Config
	data       map[string]string
}

// NewPDBCollector is a constructor
func NewPDBCollector(config *restclient.Config) *PDBCollector {
	return &PDBCollector{
		data:       make(map[string]string),
		kubeconfig: config,
	}
}

func (collector *PDBCollector) GetName() string {
	return "poddisruptionbudget"
}

// Collect implements the interface method
func (collector *PDBCollector) Collect() error {
	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(collector.kubeconfig)
	if err != nil {
		return fmt.Errorf("getting access to K8S failed: %w", err)
	}

	ctxBackground := context.Background()

	namespacesList, err := clientset.CoreV1().Namespaces().List(ctxBackground, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list namespaces in the cluster: %w", err)
	}

	for _, namespace := range namespacesList.Items {
		podDistInterface, err := clientset.PolicyV1().PodDisruptionBudgets(namespace.Name).List(ctxBackground, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("PDB error cluster: %w", err)
		}

		pdbresult := make([]PDBInfo, 0)
		for _, i := range podDistInterface.Items {
			pdbinfo := PDBInfo{
				Name:               i.Name,
				MinAvailable:       i.Spec.MinAvailable.String(),
				MaxUnavailable:     i.Spec.MaxUnavailable.String(),
				DisruptionsAllowed: i.Status.DisruptionsAllowed,
			}
			pdbresult = append(pdbresult, pdbinfo)
		}

		data, err := json.Marshal(pdbresult)

		if err != nil {
			return fmt.Errorf("marshall PDB to json: %w", err)
		}
		collector.data["pdb-"+namespace.Name] = string(data)
	}

	return nil
}

func (collector *PDBCollector) GetData() map[string]string {
	return collector.data
}
