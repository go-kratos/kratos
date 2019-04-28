package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "logging_test.go",
        "server_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/ecode/pb:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/rpc/warden/proto/testproto:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "logging.go",
        "recovery.go",
        "server.go",
        "stats.go",
        "validate.go",
    ],
    importpath = "go-common/library/net/rpc/warden",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/dsn:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/conf/flagvar:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/rpc/warden/balancer/wrr:go_default_library",
        "//library/net/rpc/warden/encoding/json:go_default_library",
        "//library/net/rpc/warden/resolver:go_default_library",
        "//library/net/rpc/warden/status:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/stat/summary:go_default_library",
        "//library/stat/sys/cpu:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/gopkg.in/go-playground/validator.v9:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//keepalive:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//peer:go_default_library",
        "@org_golang_google_grpc//reflection:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["exapmle_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/log:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/rpc/warden/proto/testproto:go_default_library",
        "//library/time:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
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
        "//library/net/rpc/warden/balancer/wrr:all-srcs",
        "//library/net/rpc/warden/benchmark/bench/client:all-srcs",
        "//library/net/rpc/warden/benchmark/bench/proto:all-srcs",
        "//library/net/rpc/warden/benchmark/bench/server:all-srcs",
        "//library/net/rpc/warden/benchmark/helloworld/client:all-srcs",
        "//library/net/rpc/warden/benchmark/helloworld/server:all-srcs",
        "//library/net/rpc/warden/encoding/json:all-srcs",
        "//library/net/rpc/warden/examples/client:all-srcs",
        "//library/net/rpc/warden/examples/grpcDebug:all-srcs",
        "//library/net/rpc/warden/examples/server:all-srcs",
        "//library/net/rpc/warden/metadata:all-srcs",
        "//library/net/rpc/warden/proto/testproto:all-srcs",
        "//library/net/rpc/warden/resolver:all-srcs",
        "//library/net/rpc/warden/status:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
