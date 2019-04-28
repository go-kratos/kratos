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
        "//app/job/main/broadcast/cmd:all-srcs",
        "//app/job/main/broadcast/conf:all-srcs",
        "//app/job/main/broadcast/server/http:all-srcs",
        "//app/job/main/broadcast/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
