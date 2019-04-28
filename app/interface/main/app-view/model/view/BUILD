package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bnj.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/app-view/model/view",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-view/model:go_default_library",
        "//app/interface/main/app-view/model/act:go_default_library",
        "//app/interface/main/app-view/model/ad:go_default_library",
        "//app/interface/main/app-view/model/bangumi:go_default_library",
        "//app/interface/main/app-view/model/creative:go_default_library",
        "//app/interface/main/app-view/model/elec:go_default_library",
        "//app/interface/main/app-view/model/game:go_default_library",
        "//app/interface/main/app-view/model/live:go_default_library",
        "//app/interface/main/app-view/model/manager:go_default_library",
        "//app/interface/main/app-view/model/special:go_default_library",
        "//app/interface/main/app-view/model/tag:go_default_library",
        "//app/interface/main/dm2/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/time:go_default_library",
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
