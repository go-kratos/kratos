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
        "//app/job/main/up-rating/cmd:all-srcs",
        "//app/job/main/up-rating/conf:all-srcs",
        "//app/job/main/up-rating/dao:all-srcs",
        "//app/job/main/up-rating/http:all-srcs",
        "//app/job/main/up-rating/model:all-srcs",
        "//app/job/main/up-rating/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
