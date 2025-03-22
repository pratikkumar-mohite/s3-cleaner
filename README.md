# S3 Cleaner
This project is a Go application designed to delete AWS S3 objects/buckets.  It retrieves the list of buckets from your AWS account and performs cleanup operations on specified bucket.

## Prerequisites

1. Go 20 and higher.
2. AWS account with IAM priviledges to perform S3 operations.

## Features

- List all S3 buckets in your AWS account.
- Delete S3 bucket objects with or without prefix.
- Delete S3 bucket
- Works with Versioned and Non-Versioned buckets.
- Utilize GO concurrency for delete operations.

## Build

1. Clone the repository:
    ```sh
    git clone https://github.com/pratikkumar-mohite/s3-cleaner.git
    ```
2. Navigate to the project directory:
    ```sh
    cd s3-cleaner
    ```
3. Build the application:
    ```sh
    make build
    ```
4. Move the binary to executable path
    ```sh
    mv s3-cleaner /usr/local/bin/
    ```

## Test
As of now the actual test are not there because we dont have s3 mock apis to mimic the s3 object behaviour specifically in go, this project has dependency on [S3Mock project](https://github.com/pratikkumar-mohite/S3Mock) to enable the `go test`.

There is an alternative to test the application with actual AWS S3 bucket with s3-cleaner.

1. Perform the above Build stage
2. Create test directory
    ```sh
    mkdir -p test/files
    ```
3. Create dummy files
    ```sh
     dd if=/dev/urandom of=test/files/file1.txt count=100 bs=1M
     dd if=/dev/urandom of=test/files/file2.txt count=100 bs=1M
    ```
4. Export AWS_UPLOAD_TEST_FILES
    ```sh
    export AWS_UPLOAD_TEST_FILES=true
    ```
5. Run the application
    ```sh
    make run
    ```
This will upload `file1.txt` and `file2.txt` to S3 bucket and then perform S3 object + bucket cleanup.


## Usage

1. Ensure you have AWS credentials configured. You can set them up using the AWS CLI (ignore if already set):
    ```sh
    aws configure --profile <your-aws-profile>
    ```
2. Run the s3-cleaner cli (Following parameters are mandatory)
    ```sh
    s3-cleaner -p pratikkumar-mohite-aws -r us-east-1 -b pratikkumar-mohite-test
    ```
3. Alternatively, Setup Environment variables and then run s3-cleaner cli
    ```sh
    export AWS_REGION=us-east-1
    export AWS_S3_BUCKET=pratikkumar-mohite-test
    export AWS_PROFILE=pratikkumar-mohite-aws
    $ s3-cleaner
    ```
4. Optional
    a. Use Prefix, In case you want to delete specific folder
    - Use CLI - `s3-cleaner -p pratikkumar-mohite-aws -r us-east-1 -b pratikkumar-mohite-test -f /prefix/path`
    - Use ENV variable - `export AWS_S3_PREFIX=/prefix/path`
    b. Use BucketDelete - If you want to delete the bucket
    - Use CLI - `s3-cleaner -p pratikkumar-mohite-aws -r us-east-1 -b pratikkumar-mohite-test --bucket-delete`
    - Use ENV variable - `export AWS_S3_BUCKET_DELETE=true`

![Usage](docs/gif/s3-cleaner-usage.gif)

## Note
1. In case you get following issue, there might problem with bucket region.
```sh
FATA[0002] Unable to list objects, %!v(MISSING)operation error S3: ListObjectsV2, https response error StatusCode: 301, RequestID: FJFNV6SB70432CZT, HostID: testB9w==, api error PermanentRedirect: The bucket you are attempting to access must be addressed using the specified endpoint. Please send all future requests to this endpoint.
```
2. Known Issue, the ListObjectsV2Input for listing the s3 object supports only 1,000 objects to list and with s3-cleaner I've increased it to 1,00,000. In case the cli failed with following message, *RETRY*!!
```sh
FATA[0296] Unable to delete bucket, %!v(MISSING)operation error S3: DeleteBucket, https response error StatusCode: 409, RequestID: 1S2TQ1F50737F6VA, HostID: zUcZVNGhxQtg5EepWlToEuKAEQwsvc7ZQnQn7y7DmhaqOJBiF5EdlJCHGbKxt1mASDD/yukxc+8hLU8cGae4PQyZWICH/nDOCIkKX2aNZ8k=, api error BucketNotEmpty: The bucket you tried to delete is not empty
```

## Contributing

Contributions are greatly appreciated. We actively manage the issues list, and try to highlight issues suitable for newcomers. The project follows the typical GitHub pull request model. Before starting any work, please either comment on an existing issue, or file a new one.