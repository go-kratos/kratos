package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "av_charge.go",
        "av_charge_statis.go",
        "bgm.go",
        "charge_ratio.go",
        "column_charge.go",
        "income.go",
        "income_daily_statis.go",
        "up_archive_statis.go",
        "up_charge.go",
        "up_info_video.go",
    ],
    importpath = "go-common/app/job/main/growup/model/income",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/time:go_default_library"],
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
