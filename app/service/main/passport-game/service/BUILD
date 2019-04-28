load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "cache.go",
        "md5.go",
        "passport.go",
        "passport_key.go",
        "passport_login.go",
        "passport_oauth.go",
        "passport_reg.go",
        "passport_renew_token.go",
        "region.go",
        "rsa.go",
        "rsa_time_hash.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/passport-game/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/passport-game/conf:go_default_library",
        "//app/service/main/passport-game/dao:go_default_library",
        "//app/service/main/passport-game/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "account_test.go",
        "md5_test.go",
        "passport_key_test.go",
        "passport_login_test.go",
        "passport_oauth_test.go",
        "passport_reg_test.go",
        "passport_renew_token_test.go",
        "passport_test.go",
        "region_test.go",
        "rsa_test.go",
        "rsa_time_hash_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/passport-game/conf:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
