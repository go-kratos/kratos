package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "guest_test.go",
        "item_test.go",
        "seat_test.go",
        "service_test.go",
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
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "area.go",
        "bulletin.go",
        "guest.go",
        "item.go",
        "place.go",
        "seat.go",
        "service.go",
        "venue.go",
        "version.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-item/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/openplatform/ticket-item/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-item/conf:go_default_library",
        "//app/service/openplatform/ticket-item/dao:go_default_library",
        "//app/service/openplatform/ticket-item/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/gopkg.in/go-playground/validator.v9:go_default_library",
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
