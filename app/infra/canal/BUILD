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
        "//app/infra/canal/cmd:all-srcs",
        "//app/infra/canal/conf:all-srcs",
        "//app/infra/canal/dao:all-srcs",
        "//app/infra/canal/http:all-srcs",
        "//app/infra/canal/infoc:all-srcs",
        "//app/infra/canal/model:all-srcs",
        "//app/infra/canal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
