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
        "//app/infra/canal:all-srcs",
        "//app/infra/config/cmd:all-srcs",
        "//app/infra/config/conf:all-srcs",
        "//app/infra/config/dao/v1:all-srcs",
        "//app/infra/config/dao/v2:all-srcs",
        "//app/infra/config/http:all-srcs",
        "//app/infra/config/model:all-srcs",
        "//app/infra/config/rpc/client:all-srcs",
        "//app/infra/config/rpc/server:all-srcs",
        "//app/infra/config/service/v1:all-srcs",
        "//app/infra/config/service/v2:all-srcs",
        "//app/infra/databus:all-srcs",
        "//app/infra/discovery:all-srcs",
        "//app/infra/notify:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
