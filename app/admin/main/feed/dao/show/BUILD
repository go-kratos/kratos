package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "app_active.go",
        "card_set_stars.go",
        "channel_tab.go",
        "dao.go",
        "event_topic.go",
        "search_web.go",
        "search_web_card.go",
    ],
    importpath = "go-common/app/admin/main/feed/dao/show",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/feed/conf:go_default_library",
        "//app/admin/main/feed/model/common:go_default_library",
        "//app/admin/main/feed/model/show:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
        "app_active_test.go",
        "card_set_stars_test.go",
        "channel_tab_test.go",
        "dao_test.go",
        "event_topic_test.go",
        "search_web_card_test.go",
        "search_web_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/feed/conf:go_default_library",
        "//app/admin/main/feed/model/show:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
