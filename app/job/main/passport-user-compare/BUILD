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
        "//app/job/main/passport-user-compare/cmd:all-srcs",
        "//app/job/main/passport-user-compare/conf:all-srcs",
        "//app/job/main/passport-user-compare/dao:all-srcs",
        "//app/job/main/passport-user-compare/http:all-srcs",
        "//app/job/main/passport-user-compare/model:all-srcs",
        "//app/job/main/passport-user-compare/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
