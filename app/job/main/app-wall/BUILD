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
        "//app/job/main/app-wall/cmd:all-srcs",
        "//app/job/main/app-wall/conf:all-srcs",
        "//app/job/main/app-wall/dao/offer:all-srcs",
        "//app/job/main/app-wall/dao/seq:all-srcs",
        "//app/job/main/app-wall/dao/unicom:all-srcs",
        "//app/job/main/app-wall/http:all-srcs",
        "//app/job/main/app-wall/model:all-srcs",
        "//app/job/main/app-wall/service/offer:all-srcs",
        "//app/job/main/app-wall/service/unicom:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
