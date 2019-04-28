package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "guard_test.go",
        "node_test.go",
        "nodes_test.go",
        "registry_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/discovery/conf:go_default_library",
        "//app/infra/discovery/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "guard.go",
        "node.go",
        "nodes.go",
        "registry.go",
    ],
    importpath = "go-common/app/infra/discovery/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/discovery/conf:go_default_library",
        "//app/infra/discovery/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
