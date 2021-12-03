// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builderoutput

import (
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/buildpacks/pkg/buildererror"
)

func TestFromJSON(t *testing.T) {
	serialized := `
{
	"error": {
		"buildpackId": "bad-buildpack",
		"buildpackVersion": "vbad",
		"errorType": 13,
		"canonicalCode": 13,
		"errorId": "abc123",
		"errorMessage": "error-message",
		"anotherThing": 123
	},
	"stats": [
		{
			"buildpackId": "buildpack-1",
			"buildpackVersion": "v1",
			"totalDurationMs": 100,
			"userDurationMs": 101,
			"anotherThing": "shouldn't cause a problem"
		},
		{
			"buildpackId": "buildpack-2",
			"buildpackVersion": "v2",
			"totalDurationMs": 200,
			"userDurationMs": 201
		}
	],
	"warnings": [
		"Some warning",
		"Some other warning"
	],
	"customImage": true
}
`

	got, err := FromJSON([]byte(serialized))
	if err != nil {
		t.Fatal(err)
	}

	want := BuilderOutput{
		Error: buildererror.Error{
			BuildpackID:      "bad-buildpack",
			BuildpackVersion: "vbad",
			Type:             buildererror.StatusInternal,
			Status:           buildererror.StatusInternal,
			ID:               "abc123",
			Message:          "error-message",
		},
		Stats: []BuilderStat{
			{
				BuildpackID:      "buildpack-1",
				BuildpackVersion: "v1",
				DurationMs:       100,
				UserDurationMs:   101,
			},
			{
				BuildpackID:      "buildpack-2",
				BuildpackVersion: "v2",
				DurationMs:       200,
				UserDurationMs:   201,
			},
		},
		Warnings: []string{
			"Some warning",
			"Some other warning",
		},
		CustomImage: true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("builder output parsing failed got: %v, want: %v", got, want)
	}
}

func TestJSON(t *testing.T) {
	b := BuilderOutput{Error: buildererror.Error{Status: buildererror.StatusInternal}}

	s, err := b.JSON()

	if err != nil {
		t.Fatalf("Failed to marshal %v: %v", b, err)
	}
	if !strings.Contains(string(s), `"canonicalCode":13,`) {
		t.Errorf(`Expected string '"canonicalCode":13,' not found in %s`, s)
	}
}