load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "event_topic.go",
        "service.go",
        "stars.go",
    ],
    importpath = "go-common/app/admin/main/feed/service/popular",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/feed/conf:go_default_library",
        "//app/admin/main/feed/dao/show:go_default_library",
        "//app/admin/main/feed/model/common:go_default_library",
        "//app/admin/main/feed/model/show:go_default_library",
        "//app/admin/main/feed/util:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
