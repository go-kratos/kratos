package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "business_test.go",
        "callback_test.go",
        "challenge_test.go",
        "dao_test.go",
        "event_test.go",
        "extra_test.go",
        "group_test.go",
        "memcache_test.go",
        "message_test.go",
        "redis_test.go",
        "search2_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/workflow/model:go_default_library",
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/admin/main/workflow/model/search:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "callback.go",
        "challenge.go",
        "dao.go",
        "event.go",
        "extra.go",
        "group.go",
        "memcache.go",
        "message.go",
        "redis.go",
        "search2.go",
        "tag.go",
    ],
    importpath = "go-common/app/admin/main/workflow/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/workflow/model:go_default_library",
        "//app/admin/main/workflow/model/manager:go_default_library",
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/admin/main/workflow/model/search:go_default_library",
        "//app/interface/main/credit/model:go_default_library",
        "//app/job/main/member/model/block:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
