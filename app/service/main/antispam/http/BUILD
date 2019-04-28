package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cond.go",
        "filter.go",
        "http.go",
        "keyword.go",
        "protocol.go",
        "regexp.go",
        "rule.go",
    ],
    importpath = "go-common/app/service/main/antispam/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/dao:go_default_library",
        "//app/service/main/antispam/model:go_default_library",
        "//app/service/main/antispam/service:go_default_library",
        "//app/service/main/antispam/util:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["cond_test.go"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/antispam/http:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
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
