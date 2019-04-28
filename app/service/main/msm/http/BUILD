load(
    "@io_bazel_rules_go//go:def.bzl",
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

go_library(
    name = "go_default_library",
    srcs = [
        "app_auth.go",
        "config.go",
        "credential.go",
        "ecode.go",
        "http.go",
    ],
    importpath = "go-common/app/service/main/msm/http",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/msm/conf:go_default_library",
        "//app/service/main/msm/model:go_default_library",
        "//app/service/main/msm/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
    ],
)
