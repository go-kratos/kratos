package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "ai_test.go",
        "filter_test.go",
        "infoc_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//app/service/main/filter/model/actriearea:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "admin.go",
        "ai.go",
        "article.go",
        "filter.go",
        "hit.go",
        "infoc.go",
        "key.go",
        "reply.go",
        "rubbish.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/filter/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/dao:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//app/service/main/filter/model/actriearea:go_default_library",
        "//app/service/main/filter/model/actriekey:go_default_library",
        "//app/service/main/filter/model/lrulist:go_default_library",
        "//app/service/main/filter/model/rpc:go_default_library",
        "//app/service/main/filter/service/ai:go_default_library",
        "//app/service/main/filter/service/area:go_default_library",
        "//app/service/main/filter/service/filter:go_default_library",
        "//app/service/main/filter/service/white:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//app/service/main/spy/rpc/client:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/GRbit/go-pcre:go_default_library",
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
        "//app/service/main/filter/service/ai:all-srcs",
        "//app/service/main/filter/service/area:all-srcs",
        "//app/service/main/filter/service/filter:all-srcs",
        "//app/service/main/filter/service/regexp:all-srcs",
        "//app/service/main/filter/service/white:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
