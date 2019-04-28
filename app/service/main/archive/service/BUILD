package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "region_pb_test.go",
        "region_test.go",
        "service_test.go",
        "shot_test.go",
        "stat_pb_test.go",
        "upper_pb_test.go",
        "upper_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/conf:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_pb.go",
        "region.go",
        "region_pb.go",
        "service.go",
        "shot.go",
        "stat_pb.go",
        "upper.go",
        "upper_pb.go",
        "video.go",
    ],
    importpath = "go-common/app/service/main/archive/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/conf:go_default_library",
        "//app/service/main/archive/dao/archive:go_default_library",
        "//app/service/main/archive/dao/share:go_default_library",
        "//app/service/main/archive/dao/videoshot:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
