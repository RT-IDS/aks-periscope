package collector

import (
	"os"
	"path"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func TestNewKubeObjectsCollector(t *testing.T) {
	tests := []struct {
		name          string
		want          int
		wantErr       bool
		collectorName string
	}{
		{
			name:          "get kube objects logs",
			want:          1,
			wantErr:       false,
			collectorName: "kubeobjects",
		},
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Cannot get user home dir: %v", err)
	}

	master := ""
	kubeconfig := path.Join(dirname, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		t.Fatalf("Cannot load kube config: %v", err)
	}

	c := NewKubeObjectsCollector(config)

	if err := os.Setenv("DIAGNOSTIC_KUBEOBJECTS_LIST", "kube-system/pod kube-system/service kube-system/deployment"); err != nil {
		t.Fatalf("Setenv: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Collect()

			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
			raw := c.GetData()

			if len(raw) < tt.want {
				t.Errorf("len(GetData()) = %v, want %v", len(raw), tt.want)
			}

			name := c.GetName()
			if name != tt.collectorName {
				t.Errorf("GetName()) = %v, want %v", name, tt.collectorName)
			}
		})
	}
}
