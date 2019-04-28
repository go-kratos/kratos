package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "audit_log_test.go",
        "command_test.go",
        "credit_log_test.go",
        "play_info_test.go",
        "scoresection_test.go",
        "sign_up_test.go",
        "up_base_info_test.go",
        "up_rank_test.go",
        "up_stat_history_test.go",
        "upcrmservice_test.go",
        "utils_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/model:go_default_library",
        "//app/admin/main/up/model/signmodel:go_default_library",
        "//app/admin/main/up/model/upcrmmodel:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit_log.go",
        "command.go",
        "credit_log.go",
        "play_info.go",
        "scoresection.go",
        "sign_up.go",
        "up_base_info.go",
        "up_rank.go",
        "up_stat_history.go",
        "upcrmservice.go",
        "utils.go",
    ],
    importpath = "go-common/app/admin/main/up/service/upcrmservice",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/dao/email:go_default_library",
        "//app/admin/main/up/dao/global:go_default_library",
        "//app/admin/main/up/dao/manager:go_default_library",
        "//app/admin/main/up/dao/upcrm:go_default_library",
        "//app/admin/main/up/model:go_default_library",
        "//app/admin/main/up/model/datamodel:go_default_library",
        "//app/admin/main/up/model/signmodel:go_default_library",
        "//app/admin/main/up/model/upcrmmodel:go_default_library",
        "//app/admin/main/up/service/cache:go_default_library",
        "//app/admin/main/up/service/data:go_default_library",
        "//app/admin/main/up/util:go_default_library",
        "//app/admin/main/up/util/mathutil:go_default_library",
        "//app/admin/main/up/util/timerqueue:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
