package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/ec2/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSVpcEndpointRouteTableAssociation_basic(t *testing.T) {
	resourceName := "aws_vpc_endpoint_route_table_association.test"
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointRouteTableAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointRouteTableAssociationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointRouteTableAssociationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVpcEndpointRouteTableAssociationImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSVpcEndpointRouteTableAssociation_disappears(t *testing.T) {
	resourceName := "aws_vpc_endpoint_route_table_association.test"
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointRouteTableAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointRouteTableAssociationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointRouteTableAssociationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceVPCEndpointRouteTableAssociation(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVpcEndpointRouteTableAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpc_endpoint_route_table_association" {
			continue
		}

		err := finder.VpcEndpointRouteTableAssociationExists(conn, rs.Primary.Attributes["vpc_endpoint_id"], rs.Primary.Attributes["route_table_id"])

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("VPC Endpoint Route Table Association %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckVpcEndpointRouteTableAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Endpoint Route Table Association ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		return finder.VpcEndpointRouteTableAssociationExists(conn, rs.Primary.Attributes["vpc_endpoint_id"], rs.Primary.Attributes["route_table_id"])
	}
}

func testAccAWSVpcEndpointRouteTableAssociationImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		id := fmt.Sprintf("%s/%s", rs.Primary.Attributes["vpc_endpoint_id"], rs.Primary.Attributes["route_table_id"])
		return id, nil
	}
}

func testAccVpcEndpointRouteTableAssociationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

data "aws_region" "current" {}

resource "aws_vpc_endpoint" "test" {
  vpc_id       = aws_vpc.test.id
  service_name = "com.amazonaws.${data.aws_region.current.name}.s3"

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint_route_table_association" "test" {
  vpc_endpoint_id = aws_vpc_endpoint.test.id
  route_table_id  = aws_route_table.test.id
}
`, rName)
}
