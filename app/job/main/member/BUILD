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
        "//app/job/main/member/cmd:all-srcs",
        "//app/job/main/member/conf:all-srcs",
        "//app/job/main/member/dao:all-srcs",
        "//app/job/main/member/http:all-srcs",
        "//app/job/main/member/model:all-srcs",
        "//app/job/main/member/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
