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
        "//app/service/main/usersuit/cmd:all-srcs",
        "//app/service/main/usersuit/conf:all-srcs",
        "//app/service/main/usersuit/dao/invite:all-srcs",
        "//app/service/main/usersuit/dao/medal:all-srcs",
        "//app/service/main/usersuit/dao/pendant:all-srcs",
        "//app/service/main/usersuit/http:all-srcs",
        "//app/service/main/usersuit/model:all-srcs",
        "//app/service/main/usersuit/rpc/client:all-srcs",
        "//app/service/main/usersuit/rpc/server:all-srcs",
        "//app/service/main/usersuit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
