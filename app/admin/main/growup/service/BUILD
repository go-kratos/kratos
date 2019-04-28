package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "authority_test.go",
        "blacklist_test.go",
        "cheat_test.go",
        "notice_test.go",
        "offlineactivity_test.go",
        "service_test.go",
        "tag_test.go",
        "trade_test.go",
        "up_info_test.go",
        "upload_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
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
        "business.go",
        "cheat.go",
        "format.go",
        "notice.go",
        "offlineactivity.go",
        "service.go",
        "speical_award.go",
        "tag.go",
        "trade.go",
        "up_info.go",
        "upload.go",
        "white.go",
    ],
    importpath = "go-common/app/admin/main/growup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/dao:go_default_library",
        "//app/admin/main/growup/dao/message:go_default_library",
        "//app/admin/main/growup/dao/resource:go_default_library",
        "//app/admin/main/growup/dao/shell:go_default_library",
        "//app/admin/main/growup/model:go_default_library",
        "//app/admin/main/growup/model/offlineactivity:go_default_library",
        "//app/admin/main/growup/util:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/growup/service/income:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
