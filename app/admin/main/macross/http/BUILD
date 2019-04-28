package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "local.go",
        "mail.go",
        "manager.go",
        "mng_auth.go",
        "mng_role.go",
        "mng_user.go",
        "package.go",
        "publish.go",
    ],
    importpath = "go-common/app/admin/main/macross/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/macross/conf:go_default_library",
        "//app/admin/main/macross/model/mail:go_default_library",
        "//app/admin/main/macross/model/package:go_default_library",
        "//app/admin/main/macross/model/publish:go_default_library",
        "//app/admin/main/macross/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
