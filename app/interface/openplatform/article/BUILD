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
        "//app/interface/openplatform/article/cache:all-srcs",
        "//app/interface/openplatform/article/cmd:all-srcs",
        "//app/interface/openplatform/article/conf:all-srcs",
        "//app/interface/openplatform/article/dao:all-srcs",
        "//app/interface/openplatform/article/http:all-srcs",
        "//app/interface/openplatform/article/model:all-srcs",
        "//app/interface/openplatform/article/rpc/client:all-srcs",
        "//app/interface/openplatform/article/rpc/server:all-srcs",
        "//app/interface/openplatform/article/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
