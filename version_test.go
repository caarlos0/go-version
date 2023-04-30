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

import (
	"runtime/debug"
	"testing"
	"time"
)

const art = ` _            _
| |_ ___  ___| |_
| __/ _ \/ __| __|
| ||  __/\__ \ |_
 \__\___||___/\__|
`

func TestVersionText(t *testing.T) {
	sut := GetVersionInfo(
		WithAppDetails("test", "a test description"),
		WithASCIIName(art),
		WithURL("https://carlosbecker.com"),
		WithBuiltBy("nixpkgs"),
	)
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

func TestGetGitVersion(t *testing.T) {
	t.Run("null buildinfo", func(t *testing.T) {
		if got := getGitVersion(nil); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
	t.Run("devel", func(t *testing.T) {
		if got := getGitVersion(&debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		}); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
	t.Run("empty", func(t *testing.T) {
		if got := getGitVersion(&debug.BuildInfo{}); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
	t.Run("versioned", func(t *testing.T) {
		v := "1.0.0"
		if got := getGitVersion(&debug.BuildInfo{
			Main: debug.Module{
				Version: v,
			},
		}); got != v {
			t.Fatalf("expected %q, got %q", v, got)
		}
	})
}

func TestGetDirty(t *testing.T) {
	t.Run(unknown, func(t *testing.T) {
		if got := getDirty(&debug.BuildInfo{}); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
	t.Run("dirty", func(t *testing.T) {
		if got := getDirty(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.modified",
					Value: "true",
				},
			},
		}); got != "dirty" {
			t.Fatalf("expected dirty, got %q", got)
		}
	})
	t.Run("clean", func(t *testing.T) {
		if got := getDirty(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.modified",
					Value: "false",
				},
			},
		}); got != "clean" {
			t.Fatalf("expected clean, got %q", got)
		}
	})
}

func TestGetBuildDate(t *testing.T) {
	t.Run(unknown, func(t *testing.T) {
		if got := getBuildDate(&debug.BuildInfo{}); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if got := getBuildDate(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.time",
					Value: "not a date",
				},
			},
		}); got != "" {
			t.Fatalf("expected an empty string, got %q", got)
		}
	})
	t.Run("time", func(t *testing.T) {
		now := time.Now()
		if got := getBuildDate(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.time",
					Value: now.Format("2006-01-02T15:04:05Z"),
				},
			},
		}); got != now.Format("2006-01-02T15:04:05") {
			t.Fatalf("expected %q, got %q", now, got)
		}
	})
}

func TestGetKey(t *testing.T) {
	t.Run("nil buildinfo", func(t *testing.T) {
		if got := getKey(nil, "any"); got != "" {
			t.Fatalf("expected an empty string, got %q", got)
		}
	})
	t.Run("valid", func(t *testing.T) {
		key := "key"
		expect := "value"
		if got := getKey(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   key,
					Value: expect,
				},
			},
		}, key); got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})
}

func TestFirstNonEmpty(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		expect := "aaa"
		if got := firstNonEmpty("", "", expect, ""); got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})
	t.Run("all empty", func(t *testing.T) {
		if got := firstNonEmpty("", "", ""); got != "" {
			t.Fatalf("expected an empty string, got %q", got)
		}
	})
}
