package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "aids.go",
        "audit.go",
        "bfs.go",
        "bottom.go",
        "http.go",
        "language.go",
        "local.go",
        "notice.go",
        "wall.go",
    ],
    importpath = "go-common/app/admin/main/app/http",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/app/conf:go_default_library",
        "//app/admin/main/app/model:go_default_library",
        "//app/admin/main/app/model/aids:go_default_library",
        "//app/admin/main/app/model/audit:go_default_library",
        "//app/admin/main/app/model/bottom:go_default_library",
        "//app/admin/main/app/model/language:go_default_library",
        "//app/admin/main/app/model/notice:go_default_library",
        "//app/admin/main/app/model/wall:go_default_library",
        "//app/admin/main/app/service/aids:go_default_library",
        "//app/admin/main/app/service/audit:go_default_library",
        "//app/admin/main/app/service/bfs:go_default_library",
        "//app/admin/main/app/service/bottom:go_default_library",
        "//app/admin/main/app/service/language:go_default_library",
        "//app/admin/main/app/service/notice:go_default_library",
        "//app/admin/main/app/service/ping:go_default_library",
        "//app/admin/main/app/service/wall:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
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
