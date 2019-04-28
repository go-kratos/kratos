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
        "//app/job/main/passport-user/cmd:all-srcs",
        "//app/job/main/passport-user/conf:all-srcs",
        "//app/job/main/passport-user/dao:all-srcs",
        "//app/job/main/passport-user/http:all-srcs",
        "//app/job/main/passport-user/model:all-srcs",
        "//app/job/main/passport-user/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
