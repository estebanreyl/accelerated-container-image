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

package builder

import (
	"context"
	"fmt"
	"io"
	"testing"

	testingresources "github.com/containerd/accelerated-container-image/cmd/convertor/testingresources"
	"github.com/opencontainers/go-digest"
)

func Test_overlaybd_builder_uploadBaseLayer(t *testing.T) {
	ctx := context.Background()

	testingresources.RunTestWithTempDir(t, ctx, "overlaybd-builder", func(t *testing.T, ctx context.Context, workdir string) {
		reg := testingresources.GetTestRegistry(t, ctx, testingresources.RegistryOptions{
			InmemoryOnly: true,
		})
		resolver := testingresources.GetCustomTestResolver(t, ctx, reg)
		pusher := testingresources.GetTestPusherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref)

		// Create a new overlaybd builder engine
		base := &builderEngineBase{
			workDir: workdir,
			pusher:  pusher,
		}
		e := &overlaybdBuilderEngine{
			builderEngineBase: base,
		}
		desc, err := e.uploadBaseLayer(ctx)
		if err != nil {
			t.Error(err)
		}
		testingresources.Assert(t,
			desc.Digest == testingresources.ExpectedOverlaybdBaseLayerDigest,
			fmt.Sprintf("Expected digest %s, got %s", testingresources.ExpectedOverlaybdBaseLayerDigest, desc.Digest))

		// Verify that the layer is in the registry and is the smae as expected
		fetcher := testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref)
		blob, err := fetcher.Fetch(ctx, desc)
		if err != nil {
			t.Error(err)
		}
		// read blob into byte
		data, err := io.ReadAll(blob)
		if err != nil {
			t.Error(err)
		}

		digest := digest.FromBytes(data)

		testingresources.Assert(t,
			digest == testingresources.ExpectedOverlaybdBaseLayerDigest,
			fmt.Sprintf("Expected digest %s, got %s", testingresources.ExpectedOverlaybdBaseLayerDigest, digest))
	})
}

func Test_overlaybd_builder_BuildLayer_alreadyConverted(t *testing.T) {

}

func Test_overlaybd_builder_CheckForConvertedLayer(t *testing.T) {
	/*We can quickly mock the DB to check for a layer
	There are three things to check here:
	1. Check if the layer is in the DB
		|
	is it in db?
		|
		yes ---------
		|           |
		no          |
			is it in the registry?
					|
					yes ----------------
					|           	   |
					no - U1        	   |   -- Removes entry from DB
						   can it be cross repo mounted?
						               |
									  yes ---------------- Adds entry to db
									   |
									   no
	*/
}

func Test_overlaybd_builder_StoreConvertedLayerDetails(t *testing.T) {
	// verify no db returns nil
	// verify e.db.CreateEntry works
	// need overlaybdLayers[0]
	// reposiotry
	// host
}
