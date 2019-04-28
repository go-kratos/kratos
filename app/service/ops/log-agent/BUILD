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
        "//app/service/ops/log-agent/cmd:all-srcs",
        "//app/service/ops/log-agent/conf:all-srcs",
        "//app/service/ops/log-agent/event:all-srcs",
        "//app/service/ops/log-agent/input:all-srcs",
        "//app/service/ops/log-agent/output:all-srcs",
        "//app/service/ops/log-agent/pipeline:all-srcs",
        "//app/service/ops/log-agent/pkg/bufio:all-srcs",
        "//app/service/ops/log-agent/pkg/common:all-srcs",
        "//app/service/ops/log-agent/pkg/flowmonitor:all-srcs",
        "//app/service/ops/log-agent/pkg/httpstream:all-srcs",
        "//app/service/ops/log-agent/pkg/lancermonitor:all-srcs",
        "//app/service/ops/log-agent/pkg/lancerroute:all-srcs",
        "//app/service/ops/log-agent/pkg/limit:all-srcs",
        "//app/service/ops/log-agent/processor:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
