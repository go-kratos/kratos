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
        "//app/job/main/card/cmd:all-srcs",
        "//app/job/main/card/conf:all-srcs",
        "//app/job/main/card/dao:all-srcs",
        "//app/job/main/card/model:all-srcs",
        "//app/job/main/card/server/http:all-srcs",
        "//app/job/main/card/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
