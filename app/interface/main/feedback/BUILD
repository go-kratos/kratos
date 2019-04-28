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
        "//app/interface/main/feedback/cmd:all-srcs",
        "//app/interface/main/feedback/conf:all-srcs",
        "//app/interface/main/feedback/dao:all-srcs",
        "//app/interface/main/feedback/http:all-srcs",
        "//app/interface/main/feedback/model:all-srcs",
        "//app/interface/main/feedback/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
