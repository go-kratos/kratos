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
        "//app/interface/main/web/cmd:all-srcs",
        "//app/interface/main/web/conf:all-srcs",
        "//app/interface/main/web/dao:all-srcs",
        "//app/interface/main/web/http:all-srcs",
        "//app/interface/main/web/model:all-srcs",
        "//app/interface/main/web/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
