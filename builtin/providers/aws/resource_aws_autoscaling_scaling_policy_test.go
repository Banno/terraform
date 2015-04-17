package aws

import(
    "fmt"
    //"math/rand"
    //"strings"
    "testing"
    //"time"

    "github.com/awslabs/aws-sdk-go/aws"
    "github.com/awslabs/aws-sdk-go/service/autoscaling"
    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
)

func TestAccAWSAutoscalingScalingPolicy_basic(t *testing.T) {
    var policy autoscaling.ScalingPolicy
    //var group autoscaling.AutoScalingGroup
    //var lc autoscaling.LaunchConfiguration

    resource.Test(t, resource.TestCase{
        PreCheck:       func () { testAccPreCheck(t) },
        Providers:      testAccProviders,
        CheckDestroy:   testAccCheckAWSAutoscalingScalingPolicyDestroy,
        Steps:          []resource.TestStep{
            resource.TestStep{
                Config: testAccAWSAutoscalingScalingPolicyConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckScalingPolicyExists("aws_autoscaling_scaling_policy.foobar", &policy),
                    resource.TestCheckResourceAttr("aws_autoscaling_scaling_policy.foobar", "adjustment_type", "ChangeInCapacity"),
                    resource.TestCheckResourceAttr("aws_autoscaling_scaling_policy.foobar", "cooldown", "300"),
                ),
            },
        },
    })
}

func testAccCheckScalingPolicyExists(n string, policy *autoscaling.ScalingPolicy) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        fmt.Printf("[DEBUG] alex %s", n)
        rs, ok := s.RootModule().Resources[n]
        fmt.Printf("[DEBUG] alex %#v", rs)
        if !ok {
            rs = rs
            return fmt.Errorf("Not found: %s", n)
        }

        conn := testAccProvider.Meta().(*AWSClient).autoscalingconn
        params := &autoscaling.DescribePoliciesInput{
            AutoScalingGroupName: aws.String(rs.Primary.Attributes["autoscaling_group_name"]),
            PolicyNames: []*string{aws.String(rs.Primary.ID)},
        }
        resp, err := conn.DescribePolicies(params)
        if err != nil {
            return err
        }
        if len(resp.ScalingPolicies) == 0 {
            return fmt.Errorf("ScalingPolicy not found")
        }

        *policy = *resp.ScalingPolicies[0]

        return nil
    }
}

func testAccCheckAWSAutoscalingScalingPolicyDestroy(s *terraform.State) error {
    //conn := testAccProvider.Meta().(*AWSClient).autoscalingconn

    //for _, rs := range s.RootModule().Resources {
    //    if rs.Type != "aws_autoscaling_group" {
    //        continue
    //    }

    //    describeGroups, err := conn.DescribeAutoScalingGroups(
    //        &autoscaling.DescribeAutoScalingGroupsInput{
    //            AutoScalingGroupNames: []*string{aws.String(rs.Primary.ID)},
    //        })

    //    if err == nil {
    //        if len(describeGroups.AutoScalingGroups) != 0 &&
    //            *describeGRoups.AutoScalingGroups[0].AutoScalingGroupName == rs.Primary.ID {
    //            return fmt.Errorf("AutoScaling Group still exists")
    //        }
    //    }
    //}

    return nil
}

var testAccAWSAutoscalingScalingPolicyConfig = fmt.Sprintf(`
    resource "aws_launch_configuration" "foobar" {
        name = "terraform-test-foobar5"
        image_id = "ami-21f78e11"
        instance_type = "t1.micro"
    }

    resource "aws_autoscaling_group" "foobar" {
        availability_zones = ["us-west-2a"]
        name = "terraform-test-foobar5"
        max_size = 5
        min_size = 2
        health_check_grace_period = 300
        health_check_type = "ELB"
        desired_capacity = 4
        force_delete = true
        termination_policies = ["OldestInstance"]
        launch_configuration = "${aws_launch_configuration.foobar.name}"
        tag {
            key = "Foo"
            value = "foo-bar"
            propagate_at_launch = true
        }
    }

    resource "aws_autoscaling_scaling_policy" "foobar" {
        policy_name = "foobar"
        scaling_adjustment = 4
        adjustment_type = "ChangeInCapacity"
        cooldown = 300
        autoscaling_group_name = "${aws_autoscaling_group.foobar.name}"
    }

`)
