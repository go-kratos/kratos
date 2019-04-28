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
        "//app/interface/main/history/api/grpc:all-srcs",
        "//app/interface/main/history/cmd:all-srcs",
        "//app/interface/main/history/conf:all-srcs",
        "//app/interface/main/history/dao/history:all-srcs",
        "//app/interface/main/history/dao/toview:all-srcs",
        "//app/interface/main/history/http:all-srcs",
        "//app/interface/main/history/model:all-srcs",
        "//app/interface/main/history/rpc/client:all-srcs",
        "//app/interface/main/history/server/gorpc:all-srcs",
        "//app/interface/main/history/server/grpc:all-srcs",
        "//app/interface/main/history/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
