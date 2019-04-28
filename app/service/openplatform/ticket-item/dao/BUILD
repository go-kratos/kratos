package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "bulletin_test.go",
        "dao_test.go",
        "guest_test.go",
        "item_test.go",
        "promotion_test.go",
        "screen_test.go",
        "seat_test.go",
        "stock_test.go",
        "tag_test.go",
        "ticket_test.go",
        "venue_test.go",
        "version_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/openplatform/ticket-item/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-item/conf:go_default_library",
        "//app/service/openplatform/ticket-item/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "banner.go",
        "bulletin.go",
        "dao.cache.go",
        "dao.go",
        "fav.go",
        "guest.go",
        "item.go",
        "project.go",
        "promotion.go",
        "screen.go",
        "seat.go",
        "stock.go",
        "tag.go",
        "ticket.go",
        "venue.go",
        "version.go",
        "wish.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-item/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/common/openplatform/random:go_default_library",
        "//app/service/openplatform/ticket-item/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-item/conf:go_default_library",
        "//app/service/openplatform/ticket-item/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
