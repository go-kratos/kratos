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
        "//app/service/main/seq-server/cmd:all-srcs",
        "//app/service/main/seq-server/conf:all-srcs",
        "//app/service/main/seq-server/dao:all-srcs",
        "//app/service/main/seq-server/http:all-srcs",
        "//app/service/main/seq-server/model:all-srcs",
        "//app/service/main/seq-server/rpc/client:all-srcs",
        "//app/service/main/seq-server/rpc/server:all-srcs",
        "//app/service/main/seq-server/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
