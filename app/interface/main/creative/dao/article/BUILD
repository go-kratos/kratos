package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "article.go",
        "dao.go",
        "draft.go",
    ],
    importpath = "go-common/app/interface/main/creative/dao/article",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/article:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "article_test.go",
        "dao_test.go",
        "draft_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/article:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
