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
        "//app/interface/live/app-room/api/http/v1:all-srcs",
        "//app/interface/live/app-room/cmd:all-srcs",
        "//app/interface/live/app-room/conf:all-srcs",
        "//app/interface/live/app-room/dao:all-srcs",
        "//app/interface/live/app-room/model:all-srcs",
        "//app/interface/live/app-room/server/http:all-srcs",
        "//app/interface/live/app-room/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
