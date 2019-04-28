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
        "//app/job/main/growup/cmd:all-srcs",
        "//app/job/main/growup/conf:all-srcs",
        "//app/job/main/growup/dao:all-srcs",
        "//app/job/main/growup/http:all-srcs",
        "//app/job/main/growup/model:all-srcs",
        "//app/job/main/growup/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
