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
        "//app/service/main/archive/dao/archive:all-srcs",
        "//app/service/main/archive/dao/share:all-srcs",
        "//app/service/main/archive/dao/videoshot:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
