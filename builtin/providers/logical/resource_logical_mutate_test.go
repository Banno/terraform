package logical

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/ec2"
)

func TestAccLogicalMutate_normal(t *testing.T) {
	var v ec2.Instance

	testCheck := func(*terraform.State) error {
		// just sleep here and see it in the console.
		//log.Println("======== SLEEPING =========")
		//		time.Sleep(5 * time.Second)

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"aws_instance.foo", &v),
					testCheck,
				),
			},
			resource.TestStep{
				Config: testAccInstanceMutate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"aws_instance.foo", &v),
					testCheck,
				),
			},
		},
	})
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	conn := testAWSAccProvider.Meta().(*aws.AWSClient).Ec2Conn()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_instance" {
			continue
		}

		// Try to find the resource
		resp, err := conn.Instances(
			[]string{rs.Primary.ID}, ec2.NewFilter())
		if err == nil {
			if len(resp.Reservations) > 0 {
				return fmt.Errorf("still exist.")
			}

			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(*ec2.Error)
		if !ok {
			return err
		}
		if ec2err.Code != "InvalidInstanceID.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckInstanceExists(n string, i *ec2.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAWSAccProvider.Meta().(*aws.AWSClient).Ec2Conn()
		resp, err := conn.Instances(
			[]string{rs.Primary.ID}, ec2.NewFilter())
		if err != nil {
			return err
		}
		if len(resp.Reservations) == 0 {
			return fmt.Errorf("Instance not found")
		}

		*i = resp.Reservations[0].Instances[0]

		return nil
	}
}

const testAccInstanceConfig = `

provider "aws" {
    region = "us-east-1"
}

resource "aws_instance" "foo" {
  ami = "ami-44e9752c"
  instance_type = "c3.large"
  security_groups = ["sg-30a25e5f"]
  subnet_id = "subnet-c7387eaa"

  key_name = "knuckolls"

  tags {
    created-by = "terraform"
    Name = "logical-mutate-tester"
  }
}
`

const testAccInstanceMutate = `

provider "aws" {
    region = "us-east-1"
}

resource "aws_instance" "foo" {
  ami = "ami-44e9752c"
  instance_type = "c3.large"
  security_groups = ["sg-30a25e5f"]
  subnet_id = "subnet-c7387eaa"

  key_name = "knuckolls"

  tags {
    created-by = "terraform"
    Name = "logical-mutate-tester"
  }
}

resource "logical_mutate" "test" {
  ssh {
    ip = "$(aws_instance.foo.internal_ip)"
    username = "ubuntu"
    key_file = "~/.ssh/id_rsa"
  }

  source = ["./test_script.sh"]
  destination = "/opt/terraform/"
  
  env {
    TEST_ENV = "test_env"
  }

  exec = "test_script.sh $TEST_ENV"
}

`

/*
package logical

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

// basic 14.04 us-east
const testAccCheckLogicalMutateConfig_instance = `
resource "aws_instance" "zookeeper_aws_instance" {
  ami = "ami-44e9752c"
  instance_type = "c3.large"
  security_groups = ["sg-30a25e5f"]
  subnet_id = "subnet-c7387eaa"

  key_name = "knuckolls"
}
`

const testAccCheckLogicalMutateConfig_instance = `
// mutate the instance via the logical mutate resource.
`

func TestAccMarathonApp_basic(t *testing.T) {

	testCheckInstanceCreated := func(app *marathon.App) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			// check that instance was created.
		}
	}

	testCheckCreate := func(app *marathon.App) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			// test that instance was mutated.
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogicalMutateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLogicalMutateConfig_instance,
				Check:  resource.ComposeTestCheckFunc(

				// pass that instance to the test
				),
			},
			resource.TestStep{
				Config: testAccCheckMarathonAppConfig_mutate,
				Check:  resource.ComposeTestCheckFunc(
				// read the instance
				// pass that instance to the test for verification.
				),
			},
		},
	})
}

func testAccReadInstance(name string, app *marathon.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("marathon_app resource not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("marathon_app resource id not set correctly: %s", name)
		}

		//log.Printf("=== testAccContainerExists: rs ===\n%#v\n", rs)

		client := testAccProvider.Meta().(*marathon.Client)

		appRead, _ := client.AppRead(rs.Primary.Attributes["name"])

		//		log.Printf("=== testAccContainerExists: appRead ===\n%#v\n", appRead)

		time.Sleep(5000 * time.Millisecond)

		*app = *appRead

		return nil
	}
}


func testAccCheckMarathonAppDestroy(s *terraform.State) error {
	// make sure ec2 instance is destroyed.
}
*/
