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
        "//app/infra/databus/cmd:all-srcs",
        "//app/infra/databus/conf:all-srcs",
        "//app/infra/databus/dao:all-srcs",
        "//app/infra/databus/dsn:all-srcs",
        "//app/infra/databus/http:all-srcs",
        "//app/infra/databus/model:all-srcs",
        "//app/infra/databus/service:all-srcs",
        "//app/infra/databus/tcp:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
