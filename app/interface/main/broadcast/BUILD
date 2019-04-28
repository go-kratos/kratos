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
        "//app/interface/main/broadcast/api/grpc/v1:all-srcs",
        "//app/interface/main/broadcast/benchmark/broadcast:all-srcs",
        "//app/interface/main/broadcast/benchmark/client:all-srcs",
        "//app/interface/main/broadcast/benchmark/push:all-srcs",
        "//app/interface/main/broadcast/benchmark/push_room:all-srcs",
        "//app/interface/main/broadcast/benchmark/tool/client:all-srcs",
        "//app/interface/main/broadcast/benchmark/tool/push:all-srcs",
        "//app/interface/main/broadcast/cmd:all-srcs",
        "//app/interface/main/broadcast/conf:all-srcs",
        "//app/interface/main/broadcast/model:all-srcs",
        "//app/interface/main/broadcast/server:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
