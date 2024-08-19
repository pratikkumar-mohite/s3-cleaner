package main

import (
	"fmt"
	"time"

	"github.com/pratikkumar-mohite/s3-cleanup/cmd/s3-cleanup"
)

func main() {
	startTime := time.Now()
	cmd.S3Cleanup()
	elapsedTime := time.Since(startTime)
	fmt.Println("Total time taken: ", elapsedTime)
}
