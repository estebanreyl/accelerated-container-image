This PR is meant to add some unit testing for some of the core portions of the userspace conversion. I have added some notes below on the added tests and notes on the remaining tests.

## builder_engine
 - uploadManifestAndConfig -> Added test
 - getBuilderEngineBase -> Added test
 - isGzipLayer -> Added Test

## builder_utils.go
- fetch - Covered by fetchManifest and fetchConfig
- fetchManifest -> Added test
- fetchConfig -> Added test
- fetchManifestAndConfig -> Seemed unecessary
- downloadLayer -> TODO
- writeConfig -> TODO
- getFileDesc -> Added test
- uploadBlob -> Added test
- uploadBytes -> Added test
- buildArchiveFromFiles -> TODO
- addFileToArchive -> TODO

## builder.go
- build -> Added test

## overlaybd_builder.go
- uploadBaseLayer -> TODO
- checkForConvertedLayer -> TODO
- storeConvertedLayer -> TODO
I am not currenly sure how best to add the remaining functions due to their use of the overlaybd binaries. In the near term these should be covered by unit tests but leaving as a future work item.

## fastoci_builder.go
Similar problem to above. Leaving for later.

## End to End tests
- The existing CI seems like the best place to add end to end tests. TBD.

# Introduced Tools
Tests for the userspace convertor are not particularly simple to make, if only because they require a lot of setup and work on both the filesystem and remote sources. To help with this I have introduced a few tools to help with testing.

## local Remotes
This is an abstraction to interact with a local registry. The registry itself supports fetching, resolving, and pushes. See the local_registry.go file for more info. Along with this there is a mocks/registry folder which allows us to load stored images into the local registry for testing for now I've kept the added images small and restricted to hello-world but more can be added easily.

# Filesystem interactivity
test_utils.go introduces RunTestWithTempDir which is a helper emulating the snapshotter tests that allows us to run a test with a temporary directory. This is useful for testing the filesystem interactions of several of the functions.

# Other
Theres also a small bugfix to builder to remove a small contention issue found while testing.