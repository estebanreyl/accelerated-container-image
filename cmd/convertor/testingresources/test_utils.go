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
	"fmt"
	"os"
	"testing"

	"github.com/containerd/containerd/remotes"
)

func GetTestResolver(t *testing.T, ctx context.Context) remotes.Resolver {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	localRegistryPath := fmt.Sprintf("%s/../testingresources/mocks/registry", cwd)
	resolver, err := NewMockLocalResolver(ctx, localRegistryPath)
	if err != nil {
		t.Error(err)
	}
	return resolver
}

func GetTestFetcherFromResolver(t *testing.T, ctx context.Context, resolver remotes.Resolver, ref string) remotes.Fetcher {
	fetcher, err := resolver.Fetcher(ctx, ref)
	if err != nil {
		t.Error(err)
	}
	return fetcher
}

func Assert(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg)
	}
}
