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
        "archive_test.go",
        "config_test.go",
        "dao_test.go",
        "first_pass_test.go",
        "oper_test.go",
        "staff_test.go",
        "type_test.go",
        "video_test.go",
        "video_track_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/archive/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "config.go",
        "dao.go",
        "first_pass.go",
        "oper.go",
        "staff.go",
        "type.go",
        "video.go",
        "video_track.go",
    ],
    importpath = "go-common/app/job/main/archive/dao/archive",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/archive/conf:go_default_library",
        "//app/job/main/archive/model/archive:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
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
