load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
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
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "draw.go",
        "es.go",
        "mysql.go",
        "redis.go",
    ],
    importpath = "go-common/app/admin/main/answer/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/answer/conf:go_default_library",
        "//app/admin/main/answer/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/golang/freetype:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/golang.org/x/image/font:go_default_library",
        "//vendor/golang.org/x/image/math/fixed:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "draw_test.go",
        "es_test.go",
        "mysql_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/answer/conf:go_default_library",
        "//app/admin/main/answer/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
