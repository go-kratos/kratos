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
        "//app/job/main/click/cmd:all-srcs",
        "//app/job/main/click/conf:all-srcs",
        "//app/job/main/click/dao:all-srcs",
        "//app/job/main/click/http:all-srcs",
        "//app/job/main/click/model:all-srcs",
        "//app/job/main/click/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
