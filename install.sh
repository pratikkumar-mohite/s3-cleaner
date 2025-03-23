#!/bin/sh
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

URL="https://github.com/pratikkumar-mohite/s3-cleaner/releases/latest/download/s3-cleaner-${OS}-${ARCH}.tar.gz"

echo "Downloading from $URL..."
curl -L $URL -o s3-cleaner.tar.gz
tar -xzf s3-cleaner.tar.gz
chmod +x s3-cleaner
sudo mv s3-cleaner /usr/local/bin/s3-cleaner
echo "s3-cleaner installed successfully!"
