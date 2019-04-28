package(default_visibility = ["//visibility:public"])

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
        "//app/service/main/identify-game/api/grpc/v1:all-srcs",
        "//app/service/main/identify-game/cmd:all-srcs",
        "//app/service/main/identify-game/conf:all-srcs",
        "//app/service/main/identify-game/dao:all-srcs",
        "//app/service/main/identify-game/model:all-srcs",
        "//app/service/main/identify-game/rpc/client:all-srcs",
        "//app/service/main/identify-game/rpc/server:all-srcs",
        "//app/service/main/identify-game/server/grpc:all-srcs",
        "//app/service/main/identify-game/server/http:all-srcs",
        "//app/service/main/identify-game/service:all-srcs",
    ],
    tags = ["automanaged"],
)
