# Usage: ./generate.sh
# Prerequisites: oras
# Generates simple hello-world images in ./registry folder
# This script serves as a way to regenerate the images in the registry folder if necessary
# and to document the taken steps to generate the test registry. Add more if new images are needed.

srcTag="linux"
srcRepo="hello-world"
srcImage="docker.io/library/$srcRepo:$srcTag"
srcRegistry="docker.io/library"
destFolder="./registry"
rm -rf $destFolder
mkdir $destFolder

echo "Begin image generation based on src image: $srcImage"

# Docker
oras cp --to-oci-layout $srcImage $destFolder/hello-world:docker-list
# Tag Some submanifests
oras cp --to-oci-layout --platform linux/arm64 $srcRegistry/hello-world:linux $destFolder/hello-world/:arm64
oras cp --to-oci-layout --platform linux/amd64 $srcRegistry/hello-world:linux $destFolder/hello-world/:amd64

# Add some sample converted manifests
cd ../../../../
make
sudo bin/convertor --oci-layout ./cmd/convertor/testingresources/mocks/registry/ --repository localregistry.io/hello-world --input-tag amd64 --oci --overlaybd amd64-converted
sudo bin/convertor --oci-layout ./cmd/convertor/testingresources/mocks/registry/ --repository localregistry.io/hello-world --input-tag arm64 --oci --overlaybd arm64-converted