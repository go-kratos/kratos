package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "access.go",
        "auth.go",
        "contain.go",
        "empty_arc.go",
        "service.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/tv/service/view",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/account:go_default_library",
        "//app/interface/main/tv/dao/archive:go_default_library",
        "//app/interface/main/tv/dao/cms:go_default_library",
        "//app/interface/main/tv/dao/favorite:go_default_library",
        "//app/interface/main/tv/dao/upper:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/upper:go_default_library",
        "//app/interface/main/tv/model/view:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
