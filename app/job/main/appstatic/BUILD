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
        "//app/job/main/appstatic/cmd:all-srcs",
        "//app/job/main/appstatic/conf:all-srcs",
        "//app/job/main/appstatic/dao:all-srcs",
        "//app/job/main/appstatic/http:all-srcs",
        "//app/job/main/appstatic/model:all-srcs",
        "//app/job/main/appstatic/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
