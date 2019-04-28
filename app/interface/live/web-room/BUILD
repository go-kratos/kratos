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
        "//app/interface/live/web-room/api/http/v1:all-srcs",
        "//app/interface/live/web-room/cmd:all-srcs",
        "//app/interface/live/web-room/conf:all-srcs",
        "//app/interface/live/web-room/dao:all-srcs",
        "//app/interface/live/web-room/model:all-srcs",
        "//app/interface/live/web-room/server/http:all-srcs",
        "//app/interface/live/web-room/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
