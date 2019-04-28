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
        "//app/service/main/archive/model/archive:all-srcs",
        "//app/service/main/archive/model/videoshot:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
