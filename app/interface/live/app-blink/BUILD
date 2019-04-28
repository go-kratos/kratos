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
        "//app/interface/live/app-blink/api/http/v1:all-srcs",
        "//app/interface/live/app-blink/cmd:all-srcs",
        "//app/interface/live/app-blink/conf:all-srcs",
        "//app/interface/live/app-blink/dao:all-srcs",
        "//app/interface/live/app-blink/model:all-srcs",
        "//app/interface/live/app-blink/server/http:all-srcs",
        "//app/interface/live/app-blink/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
