package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "callback_test.go",
        "dao_test.go",
        "es_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/workflow/conf:go_default_library",
        "//app/service/main/workflow/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "callback.go",
        "dao.go",
        "es.go",
    ],
    importpath = "go-common/app/service/main/workflow/dao",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/workflow/conf:go_default_library",
        "//app/service/main/workflow/model:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/workflow/dao/sobot:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
