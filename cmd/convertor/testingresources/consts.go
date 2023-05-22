/*
Copyright The Accelerated Container Image Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package testingresources

import (
	"context"
	"strings"

	"github.com/containerd/containerd/reference"
)

/*
This package provides a local implementation of a registry complete with
sample images of different types. Its built in such a way that we can add
more complex images as well as more complex tests are required in the
future. Note that the local registry is not particularly optimized or a
good model for how to implement a local registry but its convenient to utilize
with existing skopeo tooling. This is something that can be easily revised as more
complex scenarios arise. For now we are using abstractions from
https://pkg.go.dev/github.com/containers/image/v5 for the purpose of maintaining
compatibility with skopeo image downloads as a quick, easy and reproducible way of
adding and downloading images.
*/

const (
	// MINIMAL MANIFESTS (For unit testing)
	// DOCKER V2 (amd64)
	DockerV2_Manifest_Simple_Ref           = "sample.localstore.io/hello-world:amd64"
	DockerV2_Manifest_Simple_Digest        = "sha256:7e9b6e7ba2842c91cf49f3e214d04a7a496f8214356f41d81a6e6dcad11f11e3"
	DockerV2_Manifest_Simple_Config_Digest = "sha256:9c7a54a9a43cca047013b82af109fe963fde787f63f9e016fdc3384500c2823d"

	// DOCKER MANIFEST LIST
	Docker_Manifest_List_Ref    = "sample.localstore.io/hello-world:docker-list"
	Docker_Manifest_List_Digest = "sha256:726023f73a8fc5103fa6776d48090539042cb822531c6b751b1f6dd18cb5705d"

	// TODO
	// OCI
	OCI_Manifest_Simple_Ref           = "sample.localstore.io/hello-world:oci"
	OCI_Manifest_Simple_Digest        = ""
	OCI_Manifest_Simple_Config_Digest = ""

	// Index
	OCI_Manifest_Index_Ref    = "sample.localstore.io/hello-world:oci-index"
	OCI_Manifest_Index_Digest = ""
)

// ParseRef Parses a ref into its components: host, repository, tag/digest
func ParseRef(ctx context.Context, ref string) (string, string, string, error) {
	refspec, err := reference.Parse(ref)
	if err != nil {
		return "", "", "", err
	}
	host := refspec.Hostname()
	repository := strings.TrimPrefix(refspec.Locator, host+"/")
	object := refspec.Object
	return host, repository, object, nil
}
