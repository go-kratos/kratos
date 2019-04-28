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
        "//app/interface/main/answer/cmd:all-srcs",
        "//app/interface/main/answer/conf:all-srcs",
        "//app/interface/main/answer/dao:all-srcs",
        "//app/interface/main/answer/http:all-srcs",
        "//app/interface/main/answer/model:all-srcs",
        "//app/interface/main/answer/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
