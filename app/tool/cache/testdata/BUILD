load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.cache.go",
        "dao.go",
        "multi.go",
        "none.go",
        "single.go",
    ],
    importpath = "go-common/app/tool/cache/testdata",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "multi_test.go",
        "none_test.go",
        "single_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
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
