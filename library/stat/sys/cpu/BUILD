package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cgroup_test.go",
        "stat_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/stretchr/testify/assert:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cgroup.go",
        "cpu.go",
        "sysconfig_notcgo.go",
        "util.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:linux": [
            "cpu_linux.go",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "cpu_darwin.go",
        ],
        "//conditions:default": [
            "cpu_other.go",
        ],
    }),
    importpath = "go-common/library/stat/sys/cpu",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/log:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
