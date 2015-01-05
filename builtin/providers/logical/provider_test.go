package logical

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAWSAccProvider *schema.Provider
var testLogicalAccProvider *schema.Provider

func init() {
	testAWSAccProvider = aws.Provider().(*schema.Provider)
	testLogicalAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"aws":     testAWSAccProvider,
		"logical": testLogicalAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := aws.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AWS_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_SECRET_KEY"); v == "" {
		t.Fatal("AWS_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_REGION"); v == "" {
		log.Println("[INFO] Test: Using us-east-1 as test region")
		os.Setenv("AWS_REGION", "us-east-1")
	}
}
