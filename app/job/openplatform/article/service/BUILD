package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "retry_test.go",
        "search_test.go",
        "service_test.go",
        "sort_test.go",
        "stat_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/job/openplatform/article/conf:go_default_library",
        "//app/job/openplatform/article/dao:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "article.go",
        "author.go",
        "cron.go",
        "infoc.go",
        "read.go",
        "recommend.go",
        "search.go",
        "service.go",
        "setting.go",
        "sitemap.go",
        "sort.go",
        "stat.go",
        "tag.go",
    ],
    importpath = "go-common/app/job/openplatform/article/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/rpc/client:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/job/openplatform/article/conf:go_default_library",
        "//app/job/openplatform/article/dao:go_default_library",
        "//app/job/openplatform/article/model:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/jaytaylor/html2text:go_default_library",
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
