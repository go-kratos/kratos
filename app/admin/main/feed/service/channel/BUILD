load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "card_setup.go",
        "tab.go",
    ],
    importpath = "go-common/app/admin/main/feed/service/channel",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/feed/conf:go_default_library",
        "//app/admin/main/feed/dao/account:go_default_library",
        "//app/admin/main/feed/dao/archive:go_default_library",
        "//app/admin/main/feed/dao/pgc:go_default_library",
        "//app/admin/main/feed/dao/show:go_default_library",
        "//app/admin/main/feed/model/channel:go_default_library",
        "//app/admin/main/feed/model/common:go_default_library",
        "//app/admin/main/feed/model/show:go_default_library",
        "//app/admin/main/feed/util:go_default_library",
        "//library/log:go_default_library",
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
