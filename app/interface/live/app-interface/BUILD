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
        "//app/interface/live/app-interface/api/http/v1:all-srcs",
        "//app/interface/live/app-interface/api/http/v2:all-srcs",
        "//app/interface/live/app-interface/cmd:all-srcs",
        "//app/interface/live/app-interface/conf:all-srcs",
        "//app/interface/live/app-interface/dao:all-srcs",
        "//app/interface/live/app-interface/model:all-srcs",
        "//app/interface/live/app-interface/server/http:all-srcs",
        "//app/interface/live/app-interface/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
