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
        "//app/job/main/passport-encrypt/cmd:all-srcs",
        "//app/job/main/passport-encrypt/conf:all-srcs",
        "//app/job/main/passport-encrypt/dao:all-srcs",
        "//app/job/main/passport-encrypt/http:all-srcs",
        "//app/job/main/passport-encrypt/model:all-srcs",
        "//app/job/main/passport-encrypt/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
