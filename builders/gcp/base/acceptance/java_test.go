// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package acceptance

import (
	"testing"

	"github.com/GoogleCloudPlatform/buildpacks/internal/acceptance"
)

func init() {
	acceptance.DefineFlags()
}

func TestAcceptanceJava(t *testing.T) {
	imageCtx, cleanup := acceptance.ProvisionImages(t)
	t.Cleanup(cleanup)

	testCases := []acceptance.Test{
		{
			Name:            "simple Java application",
			App:             "simple",
			Env:             []string{"GOOGLE_ENTRYPOINT=java Main.java"},
			MustUse:         []string{javaRuntime, entrypoint},
			MustNotUse:      []string{javaEntrypoint},
			EnableCacheTest: true,
		},
		{
			Name: "Java runtime version respected",
			App:  "simple",
			// Checking runtime version 8 to ensure that it is not downloading latest Java 11 version.
			Path: "/version?want=8",
			Env: []string{
				"GOOGLE_ENTRYPOINT=javac Main.java; java Main",
				"GOOGLE_RUNTIME_VERSION=8",
			},
			MustUse:    []string{javaRuntime, entrypoint},
			MustNotUse: []string{javaEntrypoint},
		},
		{
			Name:       "Java 8 maven",
			App:        "hello_quarkus_maven",
			Env:        []string{"GOOGLE_RUNTIME_VERSION=8"},
			MustUse:    []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse: []string{entrypoint},
		},
		{
			Name:            "Java 11 maven",
			App:             "hello_quarkus_maven",
			Env:             []string{"GOOGLE_RUNTIME_VERSION=11"},
			MustUse:         []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse:      []string{entrypoint},
			EnableCacheTest: true,
		},
		{
			Name:            "Java 17 maven",
			App:             "hello_quarkus_maven",
			Env:             []string{"GOOGLE_RUNTIME_VERSION=17"},
			MustUse:         []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse:      []string{entrypoint},
			EnableCacheTest: true,
		},
		{
			Name:                "Java maven (Dev Mode)",
			App:                 "hello_quarkus_maven",
			Env:                 []string{"GOOGLE_DEVMODE=1"},
			FilesMustExist:      []string{"/layers/google.java.maven/m2", "/layers/google.java.maven/m2/bin/.devmode_rebuild.sh"},
			MustRebuildOnChange: "/workspace/src/main/java/hello/Hello.java",
		},
		// TODO(b/189485322): Enable the below test when exact version can be pinned.
		// {
		// 	// This is a separate test case from Dev mode above because it has a fixed runtime version.
		// 	// Its only purpose is to test that the metadata is set correctly.
		// 	Name: "Dev mode metadata",
		// 	App:  "hello_quarkus_maven",
		// 	Env:  []string{"GOOGLE_DEVMODE=1", "GOOGLE_RUNTIME_VERSION=11"},
		// 	BOM: []acceptance.BOMEntry{
		// 		{
		// 			Name: "java",
		// 			Metadata: map[string]interface{}{
		// 				"version": "11",
		// 			},
		// 		},
		// 		{
		// 			Name: "devmode",
		// 			Metadata: map[string]interface{}{
		// 				"devmode.sync": []interface{}{
		// 					map[string]interface{}{"dest": "/workspace", "src": "**/*.java"},
		// 					map[string]interface{}{"dest": "/workspace", "src": "**/*.kt"},
		// 					map[string]interface{}{"dest": "/workspace", "src": "**/*.scala"},
		// 					map[string]interface{}{"dest": "/workspace", "src": "**/*.groovy"},
		// 					map[string]interface{}{"dest": "/workspace", "src": "**/*.clj"},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			Name:       "Java 8 gradle",
			App:        "gradle_micronaut",
			Env:        []string{"GOOGLE_ENTRYPOINT=java -jar build/libs/helloworld-0.1-all.jar", "GOOGLE_RUNTIME_VERSION=8"},
			MustUse:    []string{javaGradle, javaRuntime, entrypoint},
			MustNotUse: []string{javaEntrypoint},
		},
		{
			Name:       "Java 11 gradle",
			App:        "gradle_micronaut",
			Env:        []string{"GOOGLE_ENTRYPOINT=java -jar build/libs/helloworld-0.1-all.jar"},
			MustUse:    []string{javaGradle, javaRuntime, entrypoint},
			MustNotUse: []string{javaEntrypoint},
		},
		{
			Name:       "Java 17 gradle",
			App:        "gradle_micronaut",
			Env:        []string{"GOOGLE_ENTRYPOINT=java -jar build/libs/helloworld-0.1-all.jar"},
			MustUse:    []string{javaGradle, javaRuntime, entrypoint},
			MustNotUse: []string{javaEntrypoint},
		},
		{
			Name:                "Java gradle (Dev Mode)",
			App:                 "gradle_micronaut",
			Env:                 []string{"GOOGLE_DEVMODE=1"},
			FilesMustExist:      []string{"/layers/google.java.gradle/cache", "/layers/google.java.gradle/cache/bin/.devmode_rebuild.sh"},
			MustRebuildOnChange: "/workspace/src/main/java/example/HelloController.java",
		},
		{
			Name:       "polyglot maven",
			App:        "polyglot-maven",
			MustUse:    []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse: []string{entrypoint},
		},
		{
			Name:       "Maven build args",
			App:        "maven_testing_profile",
			Env:        []string{"GOOGLE_BUILD_ARGS=--settings=maven_settings.xml -Dnative=false"},
			MustUse:    []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse: []string{entrypoint},
		},
		{
			Name:       "Gradle build args",
			App:        "gradle_test_env",
			Env:        []string{"GOOGLE_BUILD_ARGS=-Denv=test", "GOOGLE_ENTRYPOINT=java -jar build/libs/helloworld-0.1-all.jar"},
			MustUse:    []string{javaGradle, javaRuntime, entrypoint},
			MustNotUse: []string{javaEntrypoint},
		},
		{
			Name:    "Exploded Jar",
			App:     "exploded_jar",
			MustUse: []string{javaRuntime, javaExplodedJar},
		},
		{
			Name:              "Maven with source clearing",
			App:               "hello_quarkus_maven",
			Env:               []string{"GOOGLE_CLEAR_SOURCE=true"},
			MustUse:           []string{javaMaven, javaRuntime, javaEntrypoint, javaClearSource},
			MustNotUse:        []string{entrypoint},
			FilesMustExist:    []string{"/workspace/target/hello-1-runner.jar"},
			FilesMustNotExist: []string{"/workspace/src/main/java/hello/Hello.java", "/workspace/pom.xml"},
		},
		{
			Name:              "Gradle with source clearing",
			App:               "gradle_micronaut",
			Env:               []string{"GOOGLE_CLEAR_SOURCE=true", "GOOGLE_ENTRYPOINT=java -jar build/libs/helloworld-0.1-all.jar"},
			MustUse:           []string{javaGradle, javaRuntime, entrypoint, javaClearSource},
			MustNotUse:        []string{javaEntrypoint},
			FilesMustExist:    []string{"/workspace/build/libs/helloworld-0.1-all.jar"},
			FilesMustNotExist: []string{"/workspace/src/main/java/example/Application.java", "/workspace/build.gradle"},
		},
		{
			Name:            "Java 8 native extensions",
			App:             "native_extensions",
			Env:             []string{"GOOGLE_RUNTIME_VERSION=8"},
			MustUse:         []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse:      []string{entrypoint},
			EnableCacheTest: false,
		},
		{
			Name:            "Java 11 native extensions",
			App:             "native_extensions",
			Env:             []string{"GOOGLE_RUNTIME_VERSION=11"},
			MustUse:         []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse:      []string{entrypoint},
			EnableCacheTest: false,
		},
		{
			Name:            "Java 17 native extensions",
			App:             "native_extensions",
			Env:             []string{"GOOGLE_RUNTIME_VERSION=17"},
			MustUse:         []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse:      []string{entrypoint},
			EnableCacheTest: false,
		},
		{
			Name:       "Multi-module",
			App:        "multi_module",
			Env:        []string{"GOOGLE_BUILDABLE=hello"},
			MustUse:    []string{javaMaven, javaRuntime, javaEntrypoint},
			MustNotUse: []string{entrypoint},
		},
	}
	for _, tc := range testCases {
		tc := tc
		tc.FlakyBuildAttempts = 3

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			acceptance.TestApp(t, imageCtx, tc)
		})
	}
}
