package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "authority.go",
        "banner.go",
        "bgm.go",
        "blacklist.go",
        "budgetstatistics.go",
        "business.go",
        "cheat.go",
        "formatter.go",
        "notice.go",
        "speical_award.go",
        "tag.go",
        "trade.go",
        "up.go",
    ],
    importpath = "go-common/app/admin/main/growup/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/growup/util:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/growup/model/income:all-srcs",
        "//app/admin/main/growup/model/offlineactivity:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
