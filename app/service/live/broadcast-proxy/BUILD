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
        "//app/service/live/broadcast-proxy/api/v1:all-srcs",
        "//app/service/live/broadcast-proxy/cmd:all-srcs",
        "//app/service/live/broadcast-proxy/conf:all-srcs",
        "//app/service/live/broadcast-proxy/dispatch:all-srcs",
        "//app/service/live/broadcast-proxy/expr:all-srcs",
        "//app/service/live/broadcast-proxy/grocery:all-srcs",
        "//app/service/live/broadcast-proxy/server:all-srcs",
        "//app/service/live/broadcast-proxy/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
