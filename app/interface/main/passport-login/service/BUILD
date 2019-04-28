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
        "cookie.go",
        "cookie_proxy.go",
        "generate.go",
        "md5.go",
        "pwd.go",
        "rsa.go",
        "rsa_time_hash.go",
        "service.go",
        "token.go",
        "token_proxy.go",
        "user.go",
        "user_proxy.go",
        "utils.go",
    ],
    importpath = "go-common/app/interface/main/passport-login/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/passport-login/conf:go_default_library",
        "//app/interface/main/passport-login/dao:go_default_library",
        "//app/interface/main/passport-login/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "aes_test.go",
        "cookie_test.go",
        "md5_test.go",
        "pwd_test.go",
        "rsa_test.go",
        "rsa_time_hash_test.go",
        "service_test.go",
        "utils_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/passport-login/conf:go_default_library",
        "//app/interface/main/passport-login/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
