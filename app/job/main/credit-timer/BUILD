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
        "//app/job/main/credit-timer/cmd:all-srcs",
        "//app/job/main/credit-timer/conf:all-srcs",
        "//app/job/main/credit-timer/dao:all-srcs",
        "//app/job/main/credit-timer/http:all-srcs",
        "//app/job/main/credit-timer/model:all-srcs",
        "//app/job/main/credit-timer/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
