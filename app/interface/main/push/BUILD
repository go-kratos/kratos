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
        "//app/interface/main/push/cmd:all-srcs",
        "//app/interface/main/push/conf:all-srcs",
        "//app/interface/main/push/dao:all-srcs",
        "//app/interface/main/push/http:all-srcs",
        "//app/interface/main/push/model:all-srcs",
        "//app/interface/main/push/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
