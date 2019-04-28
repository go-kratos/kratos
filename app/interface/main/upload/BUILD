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
        "//app/interface/main/upload/cmd:all-srcs",
        "//app/interface/main/upload/conf:all-srcs",
        "//app/interface/main/upload/dao:all-srcs",
        "//app/interface/main/upload/http:all-srcs",
        "//app/interface/main/upload/model:all-srcs",
        "//app/interface/main/upload/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
