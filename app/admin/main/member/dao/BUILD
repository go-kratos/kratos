load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "bfs_test.go",
        "dao_test.go",
        "face_history_test.go",
        "member_test.go",
        "merak_test.go",
        "message_test.go",
        "monitor_test.go",
        "official_test.go",
        "passport_test.go",
        "realname_test.go",
        "redis_test.go",
        "review_test.go",
        "search_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "bfs.go",
        "dao.go",
        "face_history.go",
        "member.go",
        "merak.go",
        "message.go",
        "monitor.go",
        "official.go",
        "passport.go",
        "realname.go",
        "realname_old.go",
        "redis.go",
        "review.go",
        "search.go",
    ],
    importpath = "go-common/app/admin/main/member/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/dao/block:go_default_library",
        "//app/admin/main/member/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
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
        "//app/admin/main/member/dao/block:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
