package model

import (
	"encoding/json"
	"testing"
)

func TestGitlab(t *testing.T) {
	var (
		jsonStr = `{
		"object_kind": "merge_request",
		"user": {
		  "name": "Administrator",
		  "username": "root",
		  "avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40\u0026d=identicon"
		},
		"object_attributes": {
		  "id": 99,
		  "target_branch": "master",
		  "source_branch": "ms-viewport",
		  "source_project_id": 14,
		  "author_id": 51,
		  "assignee_id": 6,
		  "title": "MS-Viewport",
		  "created_at": "2013-12-03T17:23:34Z",
		  "updated_at": "2013-12-03T17:23:34Z",
		  "st_commits": null,
		  "st_diffs": null,
		  "milestone_id": null,
		  "state": "opened",
		  "merge_status": "unchecked",
		  "target_project_id": 14,
		  "iid": 1,
		  "description": "",
		  "source":{
			"name":"Awesome Project",
			"description":"Aut reprehenderit ut est.",
			"web_url":"http://example.com/awesome_space/awesome_project",
			"avatar_url":null,
			"git_ssh_url":"git@example.com:awesome_space/awesome_project.git",
			"git_http_url":"http://example.com/awesome_space/awesome_project.git",
			"namespace":"Awesome Space",
			"visibility_level":20,
			"path_with_namespace":"awesome_space/awesome_project",
			"default_branch":"master",
			"homepage":"http://example.com/awesome_space/awesome_project",
			"url":"http://example.com/awesome_space/awesome_project.git",
			"ssh_url":"git@example.com:awesome_space/awesome_project.git",
			"http_url":"http://example.com/awesome_space/awesome_project.git"
		  },
		  "target": {
			"name":"Awesome Project",
			"description":"Aut reprehenderit ut est.",
			"web_url":"http://example.com/awesome_space/awesome_project",
			"avatar_url":null,
			"git_ssh_url":"git@example.com:awesome_space/awesome_project.git",
			"git_http_url":"http://example.com/awesome_space/awesome_project.git",
			"namespace":"Awesome Space",
			"visibility_level":20,
			"path_with_namespace":"awesome_space/awesome_project",
			"default_branch":"master",
			"homepage":"http://example.com/awesome_space/awesome_project",
			"url":"http://example.com/awesome_space/awesome_project.git",
			"ssh_url":"git@example.com:awesome_space/awesome_project.git",
			"http_url":"http://example.com/awesome_space/awesome_project.git"
		  },
		  "last_commit": {
			"id": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
			"message": "fixed readme",
			"timestamp": "2012-01-03T23:36:29+02:00",
			"url": "http://example.com/awesome_space/awesome_project/commits/da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
			"author": {
			  "name": "GitLab dev user",
			  "email": "gitlabdev@dv6700.(none)"
			}
		  },
		  "work_in_progress": false,
		  "url": "http://example.com/diaspora/merge_requests/1",
		  "action": "open",
		  "assignee": {
			"name": "User1",
			"username": "user1",
			"avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40\u0026d=identicon"
		  }
		}
	  }
	  `
		mrHook = &HookMR{}
		err    error
	)
	err = json.Unmarshal([]byte(jsonStr), mrHook)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mrHook)
	t.Log(mrHook.User)
	t.Log(mrHook.ObjectAttributes)

	var (
		commentJSONStr = `{
			"object_kind": "note",
			"user": {
				"name": "Administrator",
				"username": "root",
				"avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40\u0026d=identicon"
			},
			"project_id": 5,
			"project":{
				"name":"Gitlab Test",
				"description":"Aut reprehenderit ut est.",
				"web_url":"http://example.com/gitlab-org/gitlab-test",
				"avatar_url":null,
				"git_ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
				"git_http_url":"http://example.com/gitlab-org/gitlab-test.git",
				"namespace":"Gitlab Org",
				"visibility_level":10,
				"path_with_namespace":"gitlab-org/gitlab-test",
				"default_branch":"master",
				"homepage":"http://example.com/gitlab-org/gitlab-test",
				"url":"http://example.com/gitlab-org/gitlab-test.git",
				"ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
				"http_url":"http://example.com/gitlab-org/gitlab-test.git"
			},
			"repository":{
				"name": "Gitlab Test",
				"url": "http://localhost/gitlab-org/gitlab-test.git",
				"description": "Aut reprehenderit ut est.",
				"homepage": "http://example.com/gitlab-org/gitlab-test"
			},
			"object_attributes": {
				"id": 1244,
				"note": "This MR needs work.",
				"noteable_type": "MergeRequest",
				"author_id": 1,
				"created_at": "2015-05-17 18:21:36 UTC",
				"updated_at": "2015-05-17 18:21:36 UTC",
				"project_id": 5,
				"attachment": null,
				"line_code": null,
				"commit_id": "",
				"noteable_id": 7,
				"system": false,
				"st_diff": null,
				"url": "http://example.com/gitlab-org/gitlab-test/merge_requests/1#note_1244"
			},
			"merge_request": {
				"id": 7,
				"target_branch": "markdown",
				"source_branch": "master",
				"source_project_id": 5,
				"author_id": 8,
				"assignee_id": 28,
				"title": "Tempora et eos debitis quae laborum et.",
				"created_at": "2015-03-01 20:12:53 UTC",
				"updated_at": "2015-03-21 18:27:27 UTC",
				"milestone_id": 11,
				"state": "opened",
				"merge_status": "cannot_be_merged",
				"target_project_id": 5,
				"iid": 1,
				"description": "Et voluptas corrupti assumenda temporibus. Architecto cum animi eveniet amet asperiores. Vitae numquam voluptate est natus sit et ad id.",
				"position": 0,
				"source":{
					"name":"Gitlab Test",
					"description":"Aut reprehenderit ut est.",
					"web_url":"http://example.com/gitlab-org/gitlab-test",
					"avatar_url":null,
					"git_ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
					"git_http_url":"http://example.com/gitlab-org/gitlab-test.git",
					"namespace":"Gitlab Org",
					"visibility_level":10,
					"path_with_namespace":"gitlab-org/gitlab-test",
					"default_branch":"master",
					"homepage":"http://example.com/gitlab-org/gitlab-test",
					"url":"http://example.com/gitlab-org/gitlab-test.git",
					"ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
					"http_url":"http://example.com/gitlab-org/gitlab-test.git"
				},
				"target": {
					"name":"Gitlab Test",
					"description":"Aut reprehenderit ut est.",
					"web_url":"http://example.com/gitlab-org/gitlab-test",
					"avatar_url":null,
					"git_ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
					"git_http_url":"http://example.com/gitlab-org/gitlab-test.git",
					"namespace":"Gitlab Org",
					"visibility_level":10,
					"path_with_namespace":"gitlab-org/gitlab-test",
					"default_branch":"master",
					"homepage":"http://example.com/gitlab-org/gitlab-test",
					"url":"http://example.com/gitlab-org/gitlab-test.git",
					"ssh_url":"git@example.com:gitlab-org/gitlab-test.git",
					"http_url":"http://example.com/gitlab-org/gitlab-test.git"
				},
				"last_commit": {
					"id": "562e173be03b8ff2efb05345d12df18815438a4b",
					"message": "Merge branch 'another-branch' into 'master'\n\nCheck in this test\n",
					"timestamp": "2015-04-08T21:00:25-07:00",
					"url": "http://example.com/gitlab-org/gitlab-test/commit/562e173be03b8ff2efb05345d12df18815438a4b",
					"author": {
						"name": "John Smith",
						"email": "john@example.com"
					}
				},
				"work_in_progress": false,
				"assignee": {
					"name": "User1",
					"username": "user1",
					"avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40\u0026d=identicon"
				}
			}
		}
		`
		hookComment = &HookComment{}
	)
	err = json.Unmarshal([]byte(commentJSONStr), hookComment)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hookComment)
	t.Log(hookComment.User)
	t.Log(hookComment.Project)
	t.Log(hookComment.Repository)
	t.Log(hookComment.ObjectAttributes)
	t.Log(hookComment.MergeRequest)
	t.Log(hookComment.Commit)
}
