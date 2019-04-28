package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "service_test.go",
        "transform_test.go",
        "trie_mgr_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/util:go_default_library",
        "//app/service/main/antispam/util/trie:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cond.go",
        "filter.go",
        "keyword.go",
        "option.go",
        "pipeline.go",
        "regexp.go",
        "rule.go",
        "scheduler.go",
        "service.go",
        "service_impl.go",
        "transform.go",
        "trie_mgr.go",
    ],
    importpath = "go-common/app/service/main/antispam/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/dao:go_default_library",
        "//app/service/main/antispam/extern:go_default_library",
        "//app/service/main/antispam/model:go_default_library",
        "//app/service/main/antispam/util:go_default_library",
        "//app/service/main/antispam/util/trie:go_default_library",
        "//library/cache:go_default_library",
        "//library/log:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = [
        "keyword_test.go",
        "regexp_content_test.go",
    ],
    tags = ["automanaged"],
    deps = ["//vendor/github.com/stretchr/testify/assert:go_default_library"],
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
