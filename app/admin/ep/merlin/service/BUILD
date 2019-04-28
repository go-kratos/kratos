package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "bilihub.go",
        "cluster.go",
        "dashboard.go",
        "devicefarm.go",
        "image.go",
        "machine.go",
        "machineV2.go",
        "machine_log.go",
        "mail.go",
        "node.go",
        "service.go",
        "task.go",
        "tree.go",
        "user.go",
    ],
    importpath = "go-common/app/admin/ep/merlin/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/merlin/conf:go_default_library",
        "//app/admin/ep/merlin/dao:go_default_library",
        "//app/admin/ep/merlin/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "machine_test.go",
        "service_test.go",
        "tree_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/ep/merlin/conf:go_default_library",
        "//app/admin/ep/merlin/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)
