# lambda logger

This is a function that will listen to incoming cloudwatch log streams and send those events to humio. It contains all the terraform that is necessary to deploy the function (lambdalogger) and a test function. It also ties them together (look at tf/connections.tf) so that you can easily invoke the testfunc and see the output.

# deploying
It has two modes: prod and regular. The process works by uploading zip file and then kicking aws via the aws cli. The regular `make deploy` will upload the -latest.zip. It is good for development. The `prod_deploy` target will upload a zip file with the git SHA on the end. Then you can update TF directly.

Similarly, there are deploy_test targets. There is no hard versioning for that.

# Adding a stream
If you add a terraform stanza like that in connections.tf it start triggering this lambda.
