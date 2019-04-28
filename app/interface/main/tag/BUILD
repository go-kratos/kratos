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
        "//app/interface/main/tag/api:all-srcs",
        "//app/interface/main/tag/cmd:all-srcs",
        "//app/interface/main/tag/conf:all-srcs",
        "//app/interface/main/tag/dao:all-srcs",
        "//app/interface/main/tag/http:all-srcs",
        "//app/interface/main/tag/model:all-srcs",
        "//app/interface/main/tag/rpc/client:all-srcs",
        "//app/interface/main/tag/rpc/server:all-srcs",
        "//app/interface/main/tag/server/grpc:all-srcs",
        "//app/interface/main/tag/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
