package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/netlify/lambdalogger/version"
)

func main() {
	lambda.Start(func(ctx context.Context, input interface{}) error {
		fmt.Println("------------- start ------------")
		fmt.Println("Version", version.SHA)
		fmt.Println("Tag", version.Tag)
		fmt.Println("------------- context ------------")
		fmt.Printf("%+v\n", ctx)
		fmt.Println("------------- input ------------")
		fmt.Printf("%+v\n", input)
		fmt.Println("------------- done ------------")
		return nil
	})
}
