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
        "//app/interface/main/web-feed/cmd:all-srcs",
        "//app/interface/main/web-feed/conf:all-srcs",
        "//app/interface/main/web-feed/dao:all-srcs",
        "//app/interface/main/web-feed/http:all-srcs",
        "//app/interface/main/web-feed/model:all-srcs",
        "//app/interface/main/web-feed/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
