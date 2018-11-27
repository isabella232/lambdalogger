#
# This is a list of ARNs to listen to for logging events
# There should be one for each function that you want to 
# send logs to the logging function
#
resource "aws_cloudwatch_log_subscription_filter" "testfunc" {
  name            = "testfunc_dev"
  log_group_name  = "${aws_cloudwatch_log_group.testfunc.name}" 
  destination_arn = "${aws_lambda_function.lambdalogger.arn}"
  filter_pattern  = ""
}