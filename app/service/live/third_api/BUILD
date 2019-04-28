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
        "//app/service/live/third_api/bvc:all-srcs",
        "//app/service/live/third_api/liveBroadcast:all-srcs",
        "//app/service/live/third_api/skyhorse:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
