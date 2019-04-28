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
        "//app/infra/notify/cmd:all-srcs",
        "//app/infra/notify/conf:all-srcs",
        "//app/infra/notify/dao:all-srcs",
        "//app/infra/notify/http:all-srcs",
        "//app/infra/notify/model:all-srcs",
        "//app/infra/notify/notify:all-srcs",
        "//app/infra/notify/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
