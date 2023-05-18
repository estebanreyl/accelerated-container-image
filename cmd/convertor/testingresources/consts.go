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

const (
	// MINIMAL MANIFESTS (For unit testing)
	// OCI
	OCI_Manifest_Simple_Ref           = "sample.localstore.io/hello-world:oci"
	OCI_Manifest_Simple_Digest        = ""
	OCI_Manifest_Simple_Config_Digest = ""

	// Index
	OCI_Manifest_Index_Ref    = "sample.localstore.io/hello-world:oci-index"
	OCI_Manifest_Index_Digest = ""

	// DOCKER
	DockerV2_Manifest_Simple_Ref           = "sample.localstore.io/hello-world:docker-v2"
	DockerV2_Manifest_Simple_Digest        = ""
	DockerV2_Manifest_Simple_Config_Digest = ""

	// DOCKER MANIFEST LIST
	Docker_Manifest_List_Ref    = "sample.localstore.io/hello-world:docker-list"
	Docker_Manifest_List_Digest = ""
)
