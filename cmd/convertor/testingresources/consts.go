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
	DockerV2_Manifest_Simple_Ref            = "sample.localstore.io/hello-world:amd64"
	DockerV2_Manifest_Simple_Digest         = "sha256:e2fc4e5012d16e7fe466f5291c476431beaa1f9b90a5c2125b493ed28e2aba57"
	DockerV2_Manifest_Simple_Size           = 861
	DockerV2_Manifest_Simple_Config_Digest  = "sha256:d2c94e258dcb3c5ac2798d32e1249e42ef01cba4841c2234249495f87264ac5a"
	DockerV2_Manifest_Simple_Config_Size    = 581
	DockerV2_Manifest_Simple_Layer_0_Digest = "sha256:c1ec31eb59444d78df06a974d155e597c894ab4cda84f08294145e845394988e"
	DockerV2_Manifest_Simple_Layer_0_Size   = 2459

	// DOCKER MANIFEST LIST
	Docker_Manifest_List_Ref    = "sample.localstore.io/hello-world:docker-list"
	Docker_Manifest_List_Digest = "sha256:341ba6b5211038e9d7253263237636996d19fbf4845e78d3946a8a9f3d0f550f"
)

const (
	// OTHER CONSTS (For unit testing)
	ExpectedOverlaybdBaseLayerDigest = "sha256:a8b5fca80efae55088290f3da8110d7742de55c2a378d5ab53226a483f390e21"
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
