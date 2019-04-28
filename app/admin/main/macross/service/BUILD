package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "mail.go",
        "mng_auth.go",
        "mng_role.go",
        "mng_user.go",
        "package.go",
        "publish.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/main/macross/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/macross/conf:go_default_library",
        "//app/admin/main/macross/dao:go_default_library",
        "//app/admin/main/macross/dao/oss:go_default_library",
        "//app/admin/main/macross/model/mail:go_default_library",
        "//app/admin/main/macross/model/manager:go_default_library",
        "//app/admin/main/macross/model/package:go_default_library",
        "//app/admin/main/macross/model/publish:go_default_library",
        "//app/admin/main/macross/tools:go_default_library",
        "//library/log:go_default_library",
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
