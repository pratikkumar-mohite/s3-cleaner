# S3 Cleaner
This project is a Go application designed to delete AWS S3 buckets.  It retrieves the list of buckets from your AWS account and performs cleanup operations on specified buckets.

## Prerequisites

1. Go 20 and higher.
2. AWS account with IAM priviledges to perform S3 operations.

## Features

- List all S3 buckets in your AWS account.
- Delete specified S3 buckets and its Objects.
- Works with Versioned and Non-Versioned buckets.
- Utilize GO concurrency for delete operations.

## Build

1. Clone the repository:
    ```sh
    git clone https://github.com/pratikkumar-mohite/s2-cleanup.git
    ```
2. Navigate to the project directory:
    ```sh
    cd s3-cleanup
    ```
3. Build the application:
    ```sh
    make build
    ```
4. Move the binary to executable path
    ```sh
    mv s3-cleanup /usr/local/bin/
    ```

## Test
As of now the actual test are not there because we dont have s3 mock apis to mimic the s3 object behaviour specifically in go, this project has dependency on [S3Mock project](https://github.com/pratikkumar-mohite/S3Mock) to enable the `go test`.

## Usage

1. Ensure you have AWS credentials configured. You can set them up using the AWS CLI(ignore if already set):
    ```sh
    aws configure --profile <your-aws-profile>
    ```
2. Setup Environment variables
    ```sh
    export AWS_REGION=us-east-1
    export AWS_DELETE_S3_BUCKET=pratikkumar-mohite-test
    export AWS_PROFILE=pratikkumar-mohite-aws
    ```
2. Run the application:
    ```sh
    s3-cleanup
    ```

## Contributing

Contributions are greatly appreciated. We actively manage the issues list, and try to highlight issues suitable for newcomers. The project follows the typical GitHub pull request model. Before starting any work, please either comment on an existing issue, or file a new one.