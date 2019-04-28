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
        "archive.go",
        "archive_update.go",
        "copyright.go",
        "http.go",
        "log.go",
        "log_delete.go",
        "mng.go",
        "mng_v2.go",
        "query.go",
        "upsert.go",
    ],
    importpath = "go-common/app/admin/main/search/http",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/search/conf:go_default_library",
        "//app/admin/main/search/dao:go_default_library",
        "//app/admin/main/search/model:go_default_library",
        "//app/admin/main/search/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
    ],
)
