#!/bin/sh
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
CLI="s3-cleaner"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

URL="https://github.com/pratikkumar-mohite/$CLI/releases/latest/download/$CLI-${OS}-${ARCH}.tar.gz"

echo "Downloading from $URL..."
curl -L $URL -o $CLI.tar.gz
tar -xzf $CLI.tar.gz
chmod +x $CLI
sudo mv $CLI /usr/local/bin/$CLI
rm -rf $CLI.tar.gz
echo "$CLI installed successfully!"
