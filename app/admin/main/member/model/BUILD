load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "block.go",
        "member.go",
        "monitor.go",
        "official.go",
        "pagination.go",
        "params.go",
        "realname.go",
        "realname_old.go",
        "review.go",
        "search.go",
    ],
    importpath = "go-common/app/admin/main/member/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/model/block:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/member/model/block:all-srcs",
        "//app/admin/main/member/model/bom:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["member_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/stretchr/testify/assert:go_default_library"],
)
