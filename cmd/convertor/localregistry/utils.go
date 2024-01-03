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

package localregistry

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
