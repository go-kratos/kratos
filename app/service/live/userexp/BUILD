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
        "//app/service/live/userexp/cmd:all-srcs",
        "//app/service/live/userexp/conf:all-srcs",
        "//app/service/live/userexp/dao:all-srcs",
        "//app/service/live/userexp/http:all-srcs",
        "//app/service/live/userexp/model:all-srcs",
        "//app/service/live/userexp/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
