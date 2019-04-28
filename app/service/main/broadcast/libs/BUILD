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
        "//app/service/main/broadcast/libs/bufio:all-srcs",
        "//app/service/main/broadcast/libs/bytes:all-srcs",
        "//app/service/main/broadcast/libs/encoding/binary:all-srcs",
        "//app/service/main/broadcast/libs/time:all-srcs",
        "//app/service/main/broadcast/libs/websocket:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
