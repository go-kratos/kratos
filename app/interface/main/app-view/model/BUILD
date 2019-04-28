package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-view/model",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/main/app-view/model/act:all-srcs",
        "//app/interface/main/app-view/model/ad:all-srcs",
        "//app/interface/main/app-view/model/bangumi:all-srcs",
        "//app/interface/main/app-view/model/creative:all-srcs",
        "//app/interface/main/app-view/model/elec:all-srcs",
        "//app/interface/main/app-view/model/game:all-srcs",
        "//app/interface/main/app-view/model/live:all-srcs",
        "//app/interface/main/app-view/model/manager:all-srcs",
        "//app/interface/main/app-view/model/region:all-srcs",
        "//app/interface/main/app-view/model/report:all-srcs",
        "//app/interface/main/app-view/model/search:all-srcs",
        "//app/interface/main/app-view/model/special:all-srcs",
        "//app/interface/main/app-view/model/tag:all-srcs",
        "//app/interface/main/app-view/model/view:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
