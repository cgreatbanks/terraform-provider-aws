package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAWSGlueCrawler_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlueCrawlerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					checkGlueCatalogCrawlerExists("aws_glue_catalog_crawler.test", "test"),
					resource.TestCheckResourceAttr(
						"aws_glue_catalog_crawler.test",
						"name",
						"test",
					),
					resource.TestCheckResourceAttr(
						"aws_glue_catalog_crawler.test",
						"database_name",
						"db_name",
					),
					resource.TestCheckResourceAttr(
						"aws_glue_catalog_crawler.test",
						"role",
						"AWSGlueServiceRoleDefault",
					),
				),
			},
		},
	})
}

func TestAccAWSGlueCrawler_customCrawlers(t *testing.T) {
	const resourceName = "aws_glue_catalog_crawler.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlueCrawlerConfigCustomClassifiers,
				Check: resource.ComposeTestCheckFunc(
					checkGlueCatalogCrawlerExists("aws_glue_catalog_crawler.test", "test"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						"test",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"database_name",
						"db_name",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"role",
						"AWSGlueServiceRoleDefault",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"table_prefix",
						"table_prefix",
					),
				),
			},
		},
	})
}

func checkGlueCatalogCrawlerExists(name string, crawlerName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		glueConn := testAccProvider.Meta().(*AWSClient).glueconn
		out, err := glueConn.GetCrawler(&glue.GetCrawlerInput{
			Name: aws.String(crawlerName),
		})

		if err != nil {
			return err
		}

		if out.Crawler == nil {
			return fmt.Errorf("no Glue Crawler found")
		}

		return nil
	}
}

const testAccGlueCrawlerConfigBasic = `
	resource "aws_glue_catalog_crawler" "test" {
	  name = "test"
	  database_name = "db_name"
	  role = "${aws_iam_role.glue.name}"
	  description = "TF-test-crawler"
	  schedule="cron(0 1 * * ? *)"
	  s3_target {
		path = "s3://bucket"
	  }
	}
	
	resource "aws_iam_role_policy_attachment" "aws-glue-service-role-default-policy-attachment" {
  		policy_arn = "arn:aws:iam::aws:policy/service-role/AWSGlueServiceRole"
  		role = "${aws_iam_role.glue.name}"
	}
	
	resource "aws_iam_role" "glue" {
  		name = "AWSGlueServiceRoleDefault"
  		assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "glue.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
	}
`

const testAccGlueCrawlerConfigCustomClassifiers = `
	resource "aws_glue_catalog_crawler" "test" {
	  name = "test"
	  database_name = "db_name"
	  role = "${aws_iam_role.glue.name}"
	  classifiers = [
		"${aws_glue_classifier.test.name}"
      ]
	  s3_target {
		path = "s3://bucket"
	  }
      table_prefix = "table_prefix"
	}

	resource "aws_glue_classifier" "test" {
  	  name = "tf-example"

  	  grok_classifier {
        classification = "example"
        grok_pattern   = "example"
      }
    }
	
	resource "aws_iam_role_policy_attachment" "aws-glue-service-role-default-policy-attachment" {
  		policy_arn = "arn:aws:iam::aws:policy/service-role/AWSGlueServiceRole"
  		role = "${aws_iam_role.glue.name}"
	}
	
	resource "aws_iam_role" "glue" {
  		name = "AWSGlueServiceRoleDefault"
  		assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "glue.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
	}
`
