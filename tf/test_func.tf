resource "aws_lambda_function" "testfunc" {
  role = "${aws_iam_role.lambdalogger.arn}"

  s3_bucket     = "${data.aws_s3_bucket.zips.id}"
  s3_key        = "lambda/test_func/test_func-latest.zip"
  function_name = "test_func_dev"
  handler       = "test_func.out"
  runtime       = "go1.x"
  memory_size   = 256
  timeout       = 300

  environment {
    variables = {
      SHA           = "latest"
      LOG_LEVEL     = "debug"
    }
  }

  tags {
    env = "${terraform.workspace == "prod" ? "prod" : "dev"}"
  }
}

resource "aws_cloudwatch_log_group" "testfunc" {
  name              = "/aws/lambda/${aws_lambda_function.testfunc.function_name}"
  retention_in_days = 5
}