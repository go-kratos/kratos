package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "relations_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/relation/conf:go_default_library",
        "//app/admin/main/relation/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "relation_log.go",
        "relations.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/main/relation/service",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/relation/conf:go_default_library",
        "//app/admin/main/relation/dao:go_default_library",
        "//app/admin/main/relation/model:go_default_library",
        "//app/service/main/relation/api:go_default_library",
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
