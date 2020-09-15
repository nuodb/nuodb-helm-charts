package integration

import (
	"fmt"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/storageclass.yaml"})

	for _, obj := range testlib.SplitAndRenderStorageClass(t, output, 4) {
		if obj.Name != "local-storage" {
			assert.True(t, obj.Provisioner == expectedProvisioner)
			b, err := strconv.ParseBool(options.SetValues["storageClass.allowVolumeExpansion"])
			assert.NoError(t, err)
			assert.EqualValues(t, b, *obj.AllowVolumeExpansion)

			// Validate encrypted and iopsPerGB. Amazon-only!
			re := regexp.MustCompile("(\\w+)-storage")
			values := re.FindStringSubmatch(obj.ObjectMeta.Name)
			assert.Equal(t, 2, len(values))
			class := values[1]
			classes := []string{"fast", "manual"}
			if Contains(classes, class) {
				encKey := fmt.Sprintf("storageClass.%s.encrypted", class)
				if enc, ok := options.SetValues[encKey]; ok {
					assert.EqualValues(t, enc, obj.Parameters["encrypted"])
				}
				iopsKey := fmt.Sprintf("storageClass.%s.iopsPerGB", class)
				if iops, ok := options.SetValues[iopsKey]; ok {
					assert.EqualValues(t, iops, obj.Parameters["iopsPerGB"])
				}
			}
		} else {
			// Validate local-storage is always created
			assert.EqualValues(t, "kubernetes.io/no-provisioner", obj.Provisioner)
		}
	}
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
	options := &helm.Options{}

	// Path to the helm chart we will test
	helmChartPath := "../../stable/storage-class"

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/storageclass.yaml"})

	for _, obj := range testlib.SplitAndRenderStorageClass(t, output, 1) {
		assert.EqualValues(t, "kubernetes.io/no-provisioner", obj.Provisioner)
	}
}