Things to add unit tests for:

# Convert

# End to End
-> This involves quite a bit of mocking so I am unsure we'll get to it but in general it would involve mocks for:
- Fetcher
- Resolver
- Pusher
- Database
- Commit, Apply, Create

Its probably better to add some form of integration tests instead as this is quite a bit of requirements and its so much mocking it kind of defeats much of the purpose when integration is totally possible.

# Builder
 We could test builder with a simple mock builder, one to confirm the flow works as expected, maybe tries multiple orders for the channel completions. Not that many functions to mock and overall straightforward I think.

# builder_engine
 3 functions to test
 - uploadManifestAndConfig
 - getBuilderEngineBase
 - isGzipLayer
 All three seem straightforward but do require some level of mocking

# builder_utils -> Start with this one
- fetch -> Mock fetcher, seems simple enough
- fetchManifest -> Mock fetcher , return both list and manifest type
- fetchConfig -> Mock fetcher with returning config
- downloadLayer -> Actually creates a folder, probably not unit testable
- writeConfig -> Writes file not really unit testable
- getFileDesc -> More io, not unit testable
- uploadBlob -> Pusher, unit testable with mock pusher
- uploadBytes -> Pusher, unit testable with mock pusher
- buildArchiveFromFiles -> More io, not unit testable
- addFileToArchive -> More io, not unit testable

# fastoci_builder
- Tons of mocking to get any tests running here, probably not worth it, should be integration tested anyway without much problem with a local registry

# overlaybd_builder
- Tons of mocking to get any tests running here, probably not worth it, should be integration tested anyway without much problem with a local registry
