package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cardset.go",
        "converge.go",
        "download.go",
        "eventtopic.go",
        "follow.go",
        "operate.go",
        "special.go",
        "tab.go",
    ],
    importpath = "go-common/app/interface/main/app-card/model/card/operate",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card/banner:go_default_library",
        "//app/interface/main/app-card/model/card/cm:go_default_library",
        "//app/interface/main/app-card/model/card/rank:go_default_library",
        "//library/log:go_default_library",
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
