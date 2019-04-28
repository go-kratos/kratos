package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "engine.go",
        "log.go",
        "report.go",
    ],
    importpath = "go-common/app/admin/main/aegis/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/aegis/model/common:go_default_library",
        "//app/admin/main/aegis/model/net:go_default_library",
        "//app/admin/main/aegis/model/resource:go_default_library",
        "//app/admin/main/aegis/model/task:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
        "//app/admin/main/aegis/model/business:all-srcs",
        "//app/admin/main/aegis/model/common:all-srcs",
        "//app/admin/main/aegis/model/databus:all-srcs",
        "//app/admin/main/aegis/model/middleware:all-srcs",
        "//app/admin/main/aegis/model/monitor:all-srcs",
        "//app/admin/main/aegis/model/net:all-srcs",
        "//app/admin/main/aegis/model/resource:all-srcs",
        "//app/admin/main/aegis/model/svg:all-srcs",
        "//app/admin/main/aegis/model/task:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["model_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/aegis/model/common:go_default_library",
        "//app/admin/main/aegis/model/resource:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
