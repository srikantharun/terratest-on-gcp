//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformGcpExample(t *testing.T) {
	t.Parallel()

	//exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../examples/terraform-gcp-example", "examples/terraform-gcp-example")
        exampleDir := "../examples/terraform-gcp-example"

	// Get the Project Id to use
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Create all resources in the following zone
	zone := "europe-west2-b"

	// Give the example bucket a unique name so we can distinguish it from any other bucket in your GCP account
	expectedBucketName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// Also give the example instance a unique name
	expectedInstanceName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// website::tag::1::Configure Terraform setting path to Terraform code, bucket name, and instance name. Construct
	// the terraform options with default retryable errors to handle the most common retryable errors in terraform
	// testing.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectId,
			"zone":           zone,
			"instance_name":  expectedInstanceName,
			"bucket_name":    expectedBucketName,
		},
	}

	// website::tag::5::At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2::This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of some of the output variables
	bucketURL := terraform.Output(t, terraformOptions, "bucket_url")
	instanceName := terraform.Output(t, terraformOptions, "instance_name")

	// website::tag::3::Verify that the new bucket url matches the expected url
	expectedURL := fmt.Sprintf("\"gs://%s\"", expectedBucketName)
        actualBucket := fmt.Sprintf("%s", bucketURL)
	assert.Equal(t, expectedURL, actualBucket)

        expectedInst := fmt.Sprintf("\"%s\"", expectedInstanceName)
        assert.Equal(t, expectedInst, instanceName)

        // Verify that the Storage Bucket exists
	gcp.AssertStorageBucketExists(t, expectedBucketName)
}

