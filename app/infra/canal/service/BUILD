package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "canal.go",
        "hbase.go",
        "instance.go",
        "master.go",
        "target.go",
        "tidb_check.go",
        "tidb_data.go",
        "tidb_instance.go",
        "tidb_proc.go",
    ],
    importpath = "go-common/app/infra/canal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/canal/conf:go_default_library",
        "//app/infra/canal/dao:go_default_library",
        "//app/infra/canal/infoc:go_default_library",
        "//app/infra/canal/model:go_default_library",
        "//app/infra/canal/service/reader:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/netutil:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/juju/errors:go_default_library",
        "//vendor/github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/canal:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/client:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/replication:go_default_library",
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
        "//app/infra/canal/service/reader:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["tidb_data_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/canal/model:go_default_library",
        "//vendor/github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog:go_default_library",
    ],
)
