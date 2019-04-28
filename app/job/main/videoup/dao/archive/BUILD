package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "addit_test.go",
        "archive_biz_test.go",
        "archive_test.go",
        "arhive_history_test.go",
        "config_test.go",
        "dao_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/videoup/conf:go_default_library",
        "//app/job/main/videoup/model/archive:go_default_library",
        "//library/database/sql:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "archive_biz.go",
        "archive_history.go",
        "config.go",
        "cover.go",
        "dao.go",
        "databus.go",
        "delay.go",
        "first_pass.go",
        "forbid.go",
        "new_video.go",
        "oper.go",
        "staff.go",
        "stat.go",
        "type.go",
        "video.go",
        "video_audit.go",
        "video_cover.go",
        "video_oper.go",
        "video_shot.go",
    ],
    importpath = "go-common/app/job/main/videoup/dao/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/videoup/conf:go_default_library",
        "//app/job/main/videoup/model/archive:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
