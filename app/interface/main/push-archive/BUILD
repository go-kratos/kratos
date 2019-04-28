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
        "//app/interface/main/push-archive/cmd:all-srcs",
        "//app/interface/main/push-archive/conf:all-srcs",
        "//app/interface/main/push-archive/dao:all-srcs",
        "//app/interface/main/push-archive/http:all-srcs",
        "//app/interface/main/push-archive/model:all-srcs",
        "//app/interface/main/push-archive/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
