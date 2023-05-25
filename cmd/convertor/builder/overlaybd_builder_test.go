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
	"testing"

	"github.com/containerd/accelerated-container-image/cmd/convertor/testingresources"
)

func Test_overlaybd_builder_BuildLayer(t *testing.T) {
	ctx := context.Background()

	testingresources.RunTestWithTempDir(t, ctx, "overlaybd-builder", func(t *testing.T, ctx context.Context, workdir string) {
		// Create a new overlaybd builder engine
		base := &builderEngineBase{
			workDir: workdir,
		}

		obdEngine := NewOverlayBDBuilderEngine(base)
		obdEngine.BuildLayer(ctx, 0)
	})
}

func Test_overlaybd_builder_UploadLayer(t *testing.T) {
	// e.getLayerDir(idx) // override
	// uploadBlob is called, mostly want to verify the blob digest.
	// Lets get a converted image and verify the digest
}

func Test_overlaybd_builder_UploadImage(t *testing.T) {
	// Not sure this one needs to be tested. Calls uploadManifestAndConfig
	// We would need to populate e.overlaybdLayers which is easy enough and thats kinda it
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

func Test_overlaybd_builder_uploadBaseLayer(t *testing.T) {
	// Temp workdir
	// verify blob is equivalent to existing blob?
	// Need a converted image to verify the blob
	// add baselayer with tar header
}

func Test_overlaybd_builder_DownloadLayer(t *testing.T) {
	// set e.manifest.Layers[0]
	// g.fetLayerDir
	// probably fetcher? I dont see much of a point to verifying this
	// as its just a small wrapper around downloadLayer
}
