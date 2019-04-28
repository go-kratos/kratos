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
        "//app/interface/main/esports/cmd:all-srcs",
        "//app/interface/main/esports/conf:all-srcs",
        "//app/interface/main/esports/dao:all-srcs",
        "//app/interface/main/esports/http:all-srcs",
        "//app/interface/main/esports/model:all-srcs",
        "//app/interface/main/esports/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
