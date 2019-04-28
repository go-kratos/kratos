package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "blacklist.go",
        "budget.go",
        "charge.go",
        "cheat.go",
        "databus.go",
        "email.go",
        "income.go",
        "lottery.go",
        "patch.go",
        "tag.go",
        "up_bill.go",
        "upload.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/growup/model",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/growup/model/charge:all-srcs",
        "//app/job/main/growup/model/income:all-srcs",
        "//app/job/main/growup/model/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
