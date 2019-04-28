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
        "//app/job/main/usersuit/cmd:all-srcs",
        "//app/job/main/usersuit/conf:all-srcs",
        "//app/job/main/usersuit/dao/medal:all-srcs",
        "//app/job/main/usersuit/dao/pendant:all-srcs",
        "//app/job/main/usersuit/http:all-srcs",
        "//app/job/main/usersuit/model:all-srcs",
        "//app/job/main/usersuit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
