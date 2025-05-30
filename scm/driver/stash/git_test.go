// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits/131cb13f4aed12e725177bc4b7c28db67839bf9f").
		Reply(200).
		Type("application/json").
		File("testdata/commit.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindCommit(context.Background(), "PRJ/my-repo", "131cb13f4aed12e725177bc4b7c28db67839bf9f")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Commit)
	raw, _ := os.ReadFile("testdata/commit.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCreateRef(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		Reply(200).
		Type("application/json").
		File("testdata/create_ref.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.CreateRef(context.Background(), "PRJ/my-repo", "scm-branch", "refs/heads/main")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/create_ref.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitDeleteRef(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Delete("rest/branch-utils/latest/projects/PRJ/repos/my-repo/branches").
		Reply(204).
		Type("application/json")

	client, _ := New("http://example.com:7990")
	resp, err := client.Git.DeleteRef(context.Background(), "PRJ/my-repo", "delete")
	if err != nil {
		t.Error(err)
	}

	if resp.Status != 204 {
		t.Errorf("DeleteRef returned %v", resp.Status)
	}
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		MatchParam("filterText", "master").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindBranch(context.Background(), "PRJ/my-repo", "master")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/tags").
		MatchParam("filterText", "v1.0.0").
		Reply(200).
		Type("application/json").
		File("testdata/tag.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindTag(context.Background(), "PRJ/my-repo", "v1.0.0")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/tag.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListCommits(t *testing.T) {
	client, _ := New("http://example.com:7990")
	_, _, err := client.Git.ListCommits(context.Background(), "PRJ/my-repo", scm.CommitListOptions{Ref: "master", Page: 1, Size: 30})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListBranches(context.Background(), "PRJ/my-repo", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/branches.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/tags").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/tags.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListTags(context.Background(), "PRJ/my-repo", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/tags.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits/131cb13f4aed12e725177bc4b7c28db67839bf9f/changes").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/changes.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListChanges(context.Background(), "PRJ/my-repo", "131cb13f4aed12e725177bc4b7c28db67839bf9f", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/changes.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCompareCommits(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/compare/changes").
		MatchParam("from", "anarbitraryshabutnotatallarbitrarylength").
		MatchParam("to", "anothershathatwillgetpaddedwithdigits121").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/changes.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.CompareCommits(context.Background(), "PRJ/my-repo", "anarbitraryshabutnotatallarbitrarylength", "anothershathatwillgetpaddedwithdigits121", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/changes.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitGetDefaultBranch(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches/default").
		Reply(200).
		Type("application/json").
		File("testdata/default_branch.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.GetDefaultBranch(context.Background(), "PRJ/my-repo")
	if err != nil {
		t.Error(err)
	}

	want := &scm.Reference{}
	raw, _ := os.ReadFile("testdata/default_branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
