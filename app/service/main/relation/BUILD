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
        "//app/service/main/relation/api:all-srcs",
        "//app/service/main/relation/cmd:all-srcs",
        "//app/service/main/relation/conf:all-srcs",
        "//app/service/main/relation/dao:all-srcs",
        "//app/service/main/relation/http:all-srcs",
        "//app/service/main/relation/model:all-srcs",
        "//app/service/main/relation/rpc/client:all-srcs",
        "//app/service/main/relation/rpc/server:all-srcs",
        "//app/service/main/relation/server/grpc:all-srcs",
        "//app/service/main/relation/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
