package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "folder.go",
        "result.go",
    ],
    importpath = "go-common/app/interface/main/app-interface/model/favorite",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/model:go_default_library",
        "//app/interface/main/app-interface/model/audio:go_default_library",
        "//app/interface/main/app-interface/model/bplus:go_default_library",
        "//app/interface/main/app-interface/model/sp:go_default_library",
        "//app/interface/main/app-interface/model/topic:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
