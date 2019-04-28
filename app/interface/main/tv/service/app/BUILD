load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "homepage_test.go",
        "module_test.go",
        "pgc_cards_test.go",
        "recommend_test.go",
        "region_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "auth.go",
        "follow.go",
        "homePage.go",
        "intervention.go",
        "module.go",
        "pages.go",
        "pgc_cards.go",
        "recommend.go",
        "region.go",
        "service.go",
        "tool.go",
        "ugc_types.go",
        "zone.go",
        "zone_index.go",
    ],
    importpath = "go-common/app/interface/main/tv/service/app",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/app:go_default_library",
        "//app/interface/main/tv/dao/archive:go_default_library",
        "//app/interface/main/tv/dao/audit:go_default_library",
        "//app/interface/main/tv/dao/cms:go_default_library",
        "//app/interface/main/tv/dao/search:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/search:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
