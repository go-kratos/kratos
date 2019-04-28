package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "aes.go",
        "captcha.go",
        "dao.go",
        "email.go",
        "mid_info.go",
        "mysql.go",
        "redis.go",
        "req_rpc.go",
        "user_act_log.go",
    ],
    importpath = "go-common/app/service/main/account-recovery/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account-recovery/conf:go_default_library",
        "//app/service/main/account-recovery/dao/sqlbuilder:go_default_library",
        "//app/service/main/account-recovery/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/rpc/client:go_default_library",
        "//app/service/main/member/api:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/account-recovery/dao/sqlbuilder:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "captcha_test.go",
        "dao_test.go",
        "email_test.go",
        "mid_info_test.go",
        "mysql_test.go",
        "redis_test.go",
        "req_rpc_test.go",
        "user_act_log_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account-recovery/conf:go_default_library",
        "//app/service/main/account-recovery/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)
