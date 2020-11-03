package observe

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"observe": testAccProvider,
	}

	if cacheDir := os.Getenv("TF_ACC_TERRAFORM_CACHE_DIR"); cacheDir != "" {
		cacheTerraformBinary(cacheDir)
	}
}

func cacheTerraformBinary(cacheDir string) {
	var (
		tfPath    = os.Getenv("TF_ACC_TERRAFORM_PATH")
		tfVersion = os.Getenv("TF_ACC_TERRAFORM_VERSION")
	)

	// do not cache if an exact path was given
	if tfPath != "" {
		return
	}

	var finder tfinstall.ExecPathFinder
	if tfVersion != "" {
		finder = tfinstall.ExactVersion(tfVersion, cacheDir)
	} else {
		tfVersion = "latest"
		finder = tfinstall.LatestVersion(cacheDir, true)
	}

	path, err := tfinstall.Find(finder)
	if err != nil {
		log.Printf("[WARN] failed to cache terraform binary: %s", err)
	}

	log.Printf("[DEBUG] downloaded version %s to %s", tfVersion, path)

	if err := os.Setenv("TF_ACC_TERRAFORM_PATH", path); err != nil {
		log.Println("[WARN] failed to set TF_ACC_TERRAFORM_PATH")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{"OBSERVE_CUSTOMER", "OBSERVE_DOMAIN"}

	for _, k := range requiredEnvVars {
		if v := os.Getenv(k); v == "" {
			t.Fatalf("%s must be set for acceptance tests", k)
		}
	}
}
