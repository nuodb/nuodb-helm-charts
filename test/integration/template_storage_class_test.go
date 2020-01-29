package integration

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"gotest.tools/assert"
	storagev1 "k8s.io/api/storage/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func StorageClassTemplateE(t *testing.T, options *helm.Options, expectedProvisioner string) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/storage-class"

	// Inject random value for `allowVolumeExpansion` to test with...
	options.SetValues["storageClass.allowVolumeExpansion"] = fmt.Sprintf("%t", rand.Float32() < 0.5)

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/storageclass.yaml"})

	partCounter := 0
	storageClasses := make([]storagev1.StorageClass, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		partCounter += 1

		var sc1 storagev1.StorageClass
		helm.UnmarshalK8SYaml(t, part, &sc1)
		storageClasses = append(storageClasses, sc1)
		
		if !strings.Contains(part, "local-storage") {
			assert.Check(t, sc1.Provisioner == expectedProvisioner)
			b, err := strconv.ParseBool(options.SetValues["storageClass.allowVolumeExpansion"])
			assert.NilError(t, err)
			assert.Check(t, *sc1.AllowVolumeExpansion == b)

			// Validate encrypted and iopsPerGB. Amazon-only!
			re := regexp.MustCompile("(\\w+)-storage")
			values := re.FindStringSubmatch(sc1.ObjectMeta.Name)
			assert.Check(t, len(values) == 2)
			class := values[1]
			classes := []string{"fast", "manual"}
			if Contains(classes, class) {
				encKey := fmt.Sprintf("storageClass.%s.encrypted", class)
				if enc, ok := options.SetValues[encKey]; ok {
					assert.Check(t, sc1.Parameters["encrypted"] == enc)
				}
				iopsKey := fmt.Sprintf("storageClass.%s.iopsPerGB", class)
				if iops, ok := options.SetValues[iopsKey]; ok {
					assert.Check(t, sc1.Parameters["iopsPerGB"] == iops)
				}
			}
		} else {
			// Validate local-storage is always created
			assert.Check(t, sc1.Provisioner == "kubernetes.io/no-provisioner")
		}
	}

	assert.Equal(t, partCounter, 4)
}

func TestStorageClassTemplateAzure(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{"cloud.provider": "azure"},
	}

	expectedProvisioner := "kubernetes.io/azure-disk"

	StorageClassTemplateE(t, options, expectedProvisioner)
}

func TestStorageClassTemplateAws(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.provider":                "amazon",
			"storageClass.fast.encrypted":   fmt.Sprintf("%t", rand.Int31n(2) == 0), // b/c golang rand is awful
			"storageClass.manual.encrypted": fmt.Sprintf("%t", rand.Int31n(2) != 0), // if we spin too fast!
			"storageClass.fast.iopsPerGB":   "120",
			"storageClass.manual.iopsPerGB": "120",
		},
	}

	expectedProvisioner := "kubernetes.io/aws-ebs"

	StorageClassTemplateE(t, options, expectedProvisioner)
}

func TestStorageClassTemplateGcp(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{"cloud.provider": "google"},
	}

	expectedProvisioner := "kubernetes.io/gce-pd"

	StorageClassTemplateE(t, options, expectedProvisioner)
}

func TestStorageClassTemplateLocal(t *testing.T) {

	options := &helm.Options{
	}

	expectedProvisioner := "kubernetes.io/no-provisioner"

	// Path to the helm chart we will test
	helmChartPath := "../../stable/storage-class"

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/storageclass.yaml"})

	partCounter := 0
	storageClasses := make([]storagev1.StorageClass, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		partCounter += 1

		var sc1 storagev1.StorageClass
		helm.UnmarshalK8SYaml(t, part, &sc1)
		storageClasses = append(storageClasses, sc1)
		
		if strings.Contains(part, "local-storage") {
			assert.Check(t, sc1.Provisioner == expectedProvisioner)
		}
	}

	assert.Equal(t, partCounter, 2)
}