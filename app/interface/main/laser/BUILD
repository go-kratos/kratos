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
        "//app/interface/main/laser/cmd:all-srcs",
        "//app/interface/main/laser/conf:all-srcs",
        "//app/interface/main/laser/dao:all-srcs",
        "//app/interface/main/laser/http:all-srcs",
        "//app/interface/main/laser/model:all-srcs",
        "//app/interface/main/laser/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
