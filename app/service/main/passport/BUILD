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
        "//app/service/main/passport/cmd:all-srcs",
        "//app/service/main/passport/conf:all-srcs",
        "//app/service/main/passport/dao:all-srcs",
        "//app/service/main/passport/http:all-srcs",
        "//app/service/main/passport/model:all-srcs",
        "//app/service/main/passport/rpc/client:all-srcs",
        "//app/service/main/passport/rpc/server:all-srcs",
        "//app/service/main/passport/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
