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
        "//app/service/main/msm/cmd:all-srcs",
        "//app/service/main/msm/conf:all-srcs",
        "//app/service/main/msm/dao:all-srcs",
        "//app/service/main/msm/http:all-srcs",
        "//app/service/main/msm/model:all-srcs",
        "//app/service/main/msm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
