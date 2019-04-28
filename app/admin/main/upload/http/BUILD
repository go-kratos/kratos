package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bucekt.go",
        "dir.go",
        "file.go",
        "http.go",
        "sniff.go",
        "upload.go",
    ],
    importpath = "go-common/app/admin/main/upload/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/upload/conf:go_default_library",
        "//app/admin/main/upload/model:go_default_library",
        "//app/admin/main/upload/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
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
