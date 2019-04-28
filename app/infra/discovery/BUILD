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
        "//app/infra/discovery/cmd:all-srcs",
        "//app/infra/discovery/conf:all-srcs",
        "//app/infra/discovery/dao:all-srcs",
        "//app/infra/discovery/http:all-srcs",
        "//app/infra/discovery/model:all-srcs",
        "//app/infra/discovery/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
