package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "app_test.go",
        "build_test.go",
        "common_test.go",
        "conf_svr_test.go",
        "config_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/config/conf:go_default_library",
        "//app/admin/main/config/model:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "apm.go",
        "app.go",
        "build.go",
        "common.go",
        "conf_svr.go",
        "config.go",
        "force.go",
        "service.go",
        "tag.go",
        "team.go",
    ],
    importpath = "go-common/app/admin/main/config/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/config/conf:go_default_library",
        "//app/admin/main/config/dao:go_default_library",
        "//app/admin/main/config/model:go_default_library",
        "//app/admin/main/config/pkg/lint:go_default_library",
        "//app/infra/config/model:go_default_library",
        "//app/infra/config/rpc/client:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
