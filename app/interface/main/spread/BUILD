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
        "//app/interface/main/spread/cmd:all-srcs",
        "//app/interface/main/spread/conf:all-srcs",
        "//app/interface/main/spread/dao:all-srcs",
        "//app/interface/main/spread/http:all-srcs",
        "//app/interface/main/spread/model:all-srcs",
        "//app/interface/main/spread/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
