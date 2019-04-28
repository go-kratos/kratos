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
        "//app/interface/main/dm/cmd:all-srcs",
        "//app/interface/main/dm/conf:all-srcs",
        "//app/interface/main/dm/dao:all-srcs",
        "//app/interface/main/dm/http:all-srcs",
        "//app/interface/main/dm/model:all-srcs",
        "//app/interface/main/dm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
