package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bigdata.go",
        "bugly.go",
        "bugly2tapd.go",
        "excel.go",
        "mail.go",
        "report.go",
        "service.go",
        "tapd.go",
        "util.go",
    ],
    importpath = "go-common/app/service/ep/footman/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ep/footman/conf:go_default_library",
        "//app/service/ep/footman/dao:go_default_library",
        "//app/service/ep/footman/model:go_default_library",
        "//library/cache:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/360EntSecGroup-Skylar/excelize:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
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
