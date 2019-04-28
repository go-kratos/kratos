package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "pipeline_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/ep/saga/conf:go_default_library",
        "//app/admin/ep/saga/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/github.com/xanzy/go-gitlab:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "basic.go",
        "branch.go",
        "collect_project.go",
        "commit.go",
        "config.go",
        "data.go",
        "job.go",
        "member.go",
        "mr.go",
        "pipeline.go",
        "project.go",
        "runner.go",
        "service.go",
        "statistics.go",
        "sync.go",
        "task.go",
        "user.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/saga/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/saga/conf:go_default_library",
        "//app/admin/ep/saga/dao:go_default_library",
        "//app/admin/ep/saga/model:go_default_library",
        "//app/admin/ep/saga/service/gitlab:go_default_library",
        "//app/admin/ep/saga/service/utils:go_default_library",
        "//app/admin/ep/saga/service/wechat:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/github.com/xanzy/go-gitlab:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/admin/ep/saga/service/gitlab:all-srcs",
        "//app/admin/ep/saga/service/mail:all-srcs",
        "//app/admin/ep/saga/service/utils:all-srcs",
        "//app/admin/ep/saga/service/wechat:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
