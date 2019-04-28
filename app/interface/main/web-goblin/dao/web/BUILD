package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "channel_test.go",
        "dao_test.go",
        "pgc_test.go",
        "recruit_test.go",
        "ugc_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/model/web:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "channel.go",
        "dao.go",
        "pgc.go",
        "recruit.go",
        "ugc.go",
    ],
    importpath = "go-common/app/interface/main/web-goblin/dao/web",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/model/web:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
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
