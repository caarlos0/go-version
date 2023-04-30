/*
Copyright 2022 The Kubernetes Authors.

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

package goversion

import "testing"

const art = ` _            _
| |_ ___  ___| |_
| __/ _ \/ __| __|
| ||  __/\__ \ |_
 \__\___||___/\__|
`

func TestVersionText(t *testing.T) {
	sut := GetVersionInfo(WithAppDetails("test", "a test description"), WithASCIIName(art))
	t.Log("\n" + sut.String())
	if sut.String() == "" {
		t.Fatal("should not be empty")
	}
}

func TestVersionJSON(t *testing.T) {
	sut := GetVersionInfo()
	json, err := sut.JSONString()
	if err != nil {
		t.Fatal("expected no error, got", err)
	}

	if string(json) == "" {
		t.Fatal("should not be empty")
	}
	t.Log("\n" + string(json))
}
