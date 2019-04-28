load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

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
        "app_auth_test.go",
        "config_test.go",
        "ecode_test.go",
        "router_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/config/model:go_default_library",
        "//app/service/main/msm/conf:go_default_library",
        "//app/service/main/msm/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app_auth.go",
        "config.go",
        "ecode.go",
        "router.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/msm/service",
    tags = ["automanaged"],
    deps = [
        "//app/infra/config/model:go_default_library",
        "//app/infra/config/rpc/client:go_default_library",
        "//app/service/main/msm/conf:go_default_library",
        "//app/service/main/msm/dao:go_default_library",
        "//app/service/main/msm/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
    ],
)
