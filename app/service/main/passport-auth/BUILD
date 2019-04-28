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
        "//app/service/main/passport-auth/api/grpc/v1:all-srcs",
        "//app/service/main/passport-auth/cmd:all-srcs",
        "//app/service/main/passport-auth/conf:all-srcs",
        "//app/service/main/passport-auth/dao:all-srcs",
        "//app/service/main/passport-auth/http:all-srcs",
        "//app/service/main/passport-auth/model:all-srcs",
        "//app/service/main/passport-auth/rpc/client:all-srcs",
        "//app/service/main/passport-auth/rpc/grpc:all-srcs",
        "//app/service/main/passport-auth/rpc/server:all-srcs",
        "//app/service/main/passport-auth/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
