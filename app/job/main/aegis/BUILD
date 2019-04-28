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
        "//app/job/main/aegis/cmd:all-srcs",
        "//app/job/main/aegis/conf:all-srcs",
        "//app/job/main/aegis/dao:all-srcs",
        "//app/job/main/aegis/model:all-srcs",
        "//app/job/main/aegis/server/http:all-srcs",
        "//app/job/main/aegis/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
