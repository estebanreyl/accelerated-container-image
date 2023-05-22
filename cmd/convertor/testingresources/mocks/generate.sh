
# Usage: ./generate.sh <registry>
# Prerequisites: skopeo and skopeo login to <registry>
# Generates simple hello-world images in ./registry folder

srcTag="linux"
srcRepo="hello-world"
srcImage="docker.io/library/$srcRepo:$srcTag"
registry=$1
destFolder="./registry"

echo "Begin image generation based on src image: $srcImage"

# Docker
oras cp --to-oci-layout $srcImage $destFolder/:hello-world-docker