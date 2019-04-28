package(default_visibility = ["//visibility:public"])

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
        "//app/service/main/dapper/agent:all-srcs",
        "//app/service/main/dapper/cmd:all-srcs",
        "//app/service/main/dapper/collector:all-srcs",
        "//app/service/main/dapper/conf:all-srcs",
        "//app/service/main/dapper/dao:all-srcs",
        "//app/service/main/dapper/model:all-srcs",
        "//app/service/main/dapper/pkg/batchwrite:all-srcs",
        "//app/service/main/dapper/pkg/collect:all-srcs",
        "//app/service/main/dapper/pkg/deliver:all-srcs",
        "//app/service/main/dapper/pkg/diskqueue:all-srcs",
        "//app/service/main/dapper/pkg/pointwrite:all-srcs",
        "//app/service/main/dapper/pkg/process:all-srcs",
        "//app/service/main/dapper/pkg/util:all-srcs",
        "//app/service/main/dapper/server/udpcollect:all-srcs",
    ],
    tags = ["automanaged"],
)
