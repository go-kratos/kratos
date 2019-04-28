load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

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
        "app_auth_test.go",
        "dao_test.go",
        "mysql_test.go",
        "tree_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/msm/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app_auth.go",
        "dao.go",
        "mysql.go",
        "tree.go",
    ],
    importpath = "go-common/app/service/main/msm/dao",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/msm/conf:go_default_library",
        "//app/service/main/msm/model:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
    ],
)
