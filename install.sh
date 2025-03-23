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

echo "Downloading ${CLI} from $URL..."
if curl -L "$URL" -o "${CLI}.tar.gz"; then
    echo "Download successful."
else
    echo "Error: Failed to download ${CLI}." >&2
    exit 1
fi

# Extract and install
tar -xzf "${CLI}.tar.gz" || { echo "Error: Failed to extract ${CLI}." >&2; exit 1; }
chmod +x "${CLI}" || { echo "Error: Failed to set executable permission." >&2; exit 1; }
sudo mv "${CLI}" /usr/local/bin/"${CLI}" || { echo "Error: Failed to move ${CLI} to /usr/local/bin." >&2; exit 1; }

# Cleanup
sudo rm -rf "${CLI}.tar.gz"
echo "${CLI} installed successfully!"
