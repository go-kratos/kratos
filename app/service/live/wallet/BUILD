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
        "//app/service/live/wallet/cmd:all-srcs",
        "//app/service/live/wallet/conf:all-srcs",
        "//app/service/live/wallet/dao:all-srcs",
        "//app/service/live/wallet/http:all-srcs",
        "//app/service/live/wallet/model:all-srcs",
        "//app/service/live/wallet/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
