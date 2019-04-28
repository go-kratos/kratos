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
        "auto.go",
        "banner.go",
        "blacklist.go",
        "budgetstatistics.go",
        "charge.go",
        "cheat.go",
        "csv.go",
        "http.go",
        "income.go",
        "local.go",
        "notice.go",
        "offlineactivity.go",
        "special_award.go",
        "tag.go",
        "trade.go",
        "upload.go",
        "web.go",
        "white.go",
        "withdraw.go",
    ],
    importpath = "go-common/app/admin/main/growup/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/dao:go_default_library",
        "//app/admin/main/growup/dao/shell:go_default_library",
        "//app/admin/main/growup/model:go_default_library",
        "//app/admin/main/growup/model/offlineactivity:go_default_library",
        "//app/admin/main/growup/service:go_default_library",
        "//app/admin/main/growup/service/income:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
