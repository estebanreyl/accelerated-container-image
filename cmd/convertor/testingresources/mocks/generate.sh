
# Usage: ./generate.sh <registry>
# Prerequisites: skopeo and skopeo login to <registry>
# Generates simple hello-world images in ./registry folder

srcTag="linux"
srcRepo="hello-world"
srcImage="docker.io/library/$srcRepo:$srcTag"
registry=$1
destFolder="./registry2"

echo "Begin image generation based on src image: $srcImage"

# Docker v2
skopeo sync --src docker --dest dir $srcImage $destFolder
mv $destFolder/$srcRepo:$srcTag $destFolder/$srcRepo:docker-v2

# Docker List
skopeo sync --all --src docker --dest dir $srcImage $destFolder
mv $destFolder/$srcRepo:$srcTag $destFolder/$srcRepo:docker-list

# OCI
skopeo copy -f oci docker://$srcImage docker://$registry/$srcRepo:oci
skopeo sync --src docker --dest dir $registry/$srcRepo:oci $destFolder

# OCI Index
skopeo copy --all -f oci docker://$srcImage docker://$registry/$srcRepo:oci-index
skopeo sync --all --src docker --dest dir $registry/$srcRepo:oci-index $destFolder

v2Digest=$(skopeo manifest-digest $destFolder/$srcRepo:docker-v2/manifest.json)
v2listDigest=$(skopeo manifest-digest $destFolder/$srcRepo:docker-list/manifest.json)
ociDigest=$(skopeo manifest-digest $destFolder/$srcRepo:oci/manifest.json)
ociIndexDigest=$(skopeo manifest-digest $destFolder/$srcRepo:oci-index/manifest.json)

echo "Generated images:"
echo "docker-v2 digest: $v2Digest"
echo "docker-list digest: $v2listDigest"
echo "oci digest: $ociDigest"
echo "oci-index digest: $ociIndexDigest"