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
        "//app/interface/live/live-demo/api/http:all-srcs",
        "//app/interface/live/live-demo/cmd:all-srcs",
        "//app/interface/live/live-demo/conf:all-srcs",
        "//app/interface/live/live-demo/dao:all-srcs",
        "//app/interface/live/live-demo/server/http:all-srcs",
        "//app/interface/live/live-demo/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
