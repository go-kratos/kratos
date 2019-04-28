package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bucket.go",
        "channel.go",
        "errors.go",
        "noauth.go",
        "operation.go",
        "report.go",
        "ring.go",
        "room.go",
        "round.go",
        "server.go",
        "server_tcp.go",
        "server_tcp_v1.go",
        "server_websocket.go",
        "server_websocket_v1.go",
        "whitelist.go",
    ],
    importpath = "go-common/app/interface/main/broadcast/server",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/broadcast/conf:go_default_library",
        "//app/interface/main/broadcast/model:go_default_library",
        "//app/service/main/broadcast/api/grpc/v1:go_default_library",
        "//app/service/main/broadcast/libs/bufio:go_default_library",
        "//app/service/main/broadcast/libs/bytes:go_default_library",
        "//app/service/main/broadcast/libs/time:go_default_library",
        "//app/service/main/broadcast/libs/websocket:go_default_library",
        "//app/service/main/broadcast/model:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/netutil:go_default_library",
        "//vendor/github.com/google/uuid:go_default_library",
        "//vendor/github.com/zhenjl/cityhash:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//encoding/gzip:go_default_library",
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
        "//app/interface/main/broadcast/server/grpc:all-srcs",
        "//app/interface/main/broadcast/server/http:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["report_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
