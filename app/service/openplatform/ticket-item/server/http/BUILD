package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "area.go",
        "bulletin.go",
        "guest.go",
        "http.go",
        "item.go",
        "place.go",
        "seat.go",
        "venue.go",
        "version.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-item/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/openplatform/ticket-item/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-item/conf:go_default_library",
        "//app/service/openplatform/ticket-item/model:go_default_library",
        "//app/service/openplatform/ticket-item/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
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
