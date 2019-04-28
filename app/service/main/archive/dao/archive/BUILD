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
        "archive_mc_pb_test.go",
        "archive_test.go",
        "dao_test.go",
        "player_test.go",
        "region_rds_test.go",
        "region_test.go",
        "report_result_test.go",
        "stat_mc_pb_test.go",
        "stat_pb_test.go",
        "stat_test.go",
        "type_test.go",
        "upper_rds_test.go",
        "upper_test.go",
        "video_mc_pb_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/conf:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "archive_mc_pb.go",
        "dao.go",
        "player.go",
        "region.go",
        "region_rds.go",
        "report_result.go",
        "stat.go",
        "stat_mc_pb.go",
        "stat_pb.go",
        "type.go",
        "upper.go",
        "upper_rds.go",
        "video.go",
        "video_mc_pb.go",
    ],
    importpath = "go-common/app/service/main/archive/dao/archive",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/archive/model/databus:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/conf:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
