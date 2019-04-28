package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "credit_log_test.go",
        "dao_test.go",
        "scoresection_test.go",
        "sign_up_test.go",
        "task_info_test.go",
        "up_base_info_test.go",
        "up_play_info_test.go",
        "up_rank_test.go",
        "up_scores_history_test.go",
        "up_stats_history_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/model/signmodel:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "credit_log.go",
        "dao.go",
        "scoresection.go",
        "sign_up.go",
        "task_info.go",
        "up_base_info.go",
        "up_play_info.go",
        "up_rank.go",
        "up_scores_history.go",
        "up_stats_history.go",
    ],
    importpath = "go-common/app/admin/main/up/dao/upcrm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/model/signmodel:go_default_library",
        "//app/admin/main/up/model/upcrmmodel:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
