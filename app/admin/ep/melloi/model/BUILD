package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "apply.go",
        "bfs.go",
        "clientmoni.go",
        "cluster.go",
        "comment.go",
        "const.go",
        "dto.go",
        "grpc.go",
        "job.go",
        "label.go",
        "melloi.go",
        "order.go",
        "order_admin.go",
        "order_report.go",
        "ptestjob.go",
        "rank.go",
        "report_graph.go",
        "report_summary.go",
        "report_timely.go",
        "scene.go",
        "script.go",
        "script_snap.go",
        "tree.go",
        "user.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/melloi/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
