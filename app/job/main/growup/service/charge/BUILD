package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "av_charge_statis_test.go",
        "run_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model/charge:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "av_charge.go",
        "av_charge_statis.go",
        "bgm_charge.go",
        "column_charge.go",
        "date_statis.go",
        "run.go",
        "run_bgm.go",
        "run_column.go",
        "run_video.go",
        "service.go",
        "trans.go",
        "up_charge.go",
    ],
    importpath = "go-common/app/job/main/growup/service/charge",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/dao/charge:go_default_library",
        "//app/job/main/growup/dao/email:go_default_library",
        "//app/job/main/growup/model/charge:go_default_library",
        "//app/job/main/growup/service:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/golang.org/x/sync/errgroup:go_default_library",
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
