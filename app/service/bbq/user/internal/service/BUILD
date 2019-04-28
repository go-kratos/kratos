package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "like.go",
        "relation.go",
        "service.go",
        "user.go",
        "user_face.go",
    ],
    importpath = "go-common/app/service/bbq/user/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/user/internal/conf:go_default_library",
        "//app/service/bbq/user/internal/dao:go_default_library",
        "//app/service/bbq/user/internal/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
        "like_test.go",
        "relation_test.go",
        "service_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/user/internal/conf:go_default_library",
        "//app/service/bbq/user/internal/model:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
