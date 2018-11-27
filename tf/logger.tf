variable "project" {
  default = "lambdalogger"
}

resource "aws_lambda_function" "lambdalogger" {
  role = "${aws_iam_role.lambdalogger.arn}"

  s3_bucket     = "${data.aws_s3_bucket.zips.id}"
  s3_key        = "${var.project}/${var.project}-latest.zip"
  function_name = "${var.project}_dev"
  handler       = "${var.project}.out"
  runtime       = "go1.x"
  memory_size   = 256
  timeout       = 300

  environment {
    variables = {
      SHA           = "latest"
      LOG_LEVEL     = "debug"

      HUMIO_TOKEN      = "XKQuY2DL27gJ8CKlHy63ktinFGAgL0gLidGxjsQBHnvB"
      HUMIO_REPOSITORY = "sandbox"
    }
  }

  tags {
    env = "dev"
  }
}

# have a group for us to actually log to
resource "aws_cloudwatch_log_group" "lambdalogger" {
  name              = "/aws/lambda/${aws_lambda_function.lambdalogger.function_name}"
  retention_in_days = 5
}

# let us be triggered by writes to cloudwatch
resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  principal     = "logs.us-east-1.amazonaws.com"
  function_name = "${aws_lambda_function.lambdalogger.function_name}"
  source_arn    = "arn:aws:logs:us-east-1:${data.aws_caller_identity.current.account_id}:log-group:*:*"
}

#
# IAM
#
resource "aws_iam_role" "lambdalogger" {
    name = "lambdalogger_dev"

   assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


# this is a POLICY to attach to the above role so we can 
# send logs to cloudwatch, which is necessary for our own output
resource "aws_iam_role_policy" "cloudwatch_logs" {
    name = "cloudwatch_log_access"
    role = "${aws_iam_role.lambdalogger.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:*"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

# S3 access so that we can load our own zip file
# TODO: scope it down to read only
resource "aws_iam_role_policy" "s3_zips" {
    name = "s3_access"
    role = "${aws_iam_role.lambdalogger.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "arn:aws:s3:::${data.aws_s3_bucket.zips.id}/*"
      ]
    }
  ]
}
EOF
}

# 
# PLUMBING
#
variable profile {
  default = "personal"
}

provider "aws" {
    version = "1.14.1"
    region  = "us-east-1"
    profile = "${var.profile}"
}

 # our current account
data "aws_caller_identity" "current" {}
data "aws_s3_bucket" "zips" {
    bucket = "rybit-lambda-zips"
}