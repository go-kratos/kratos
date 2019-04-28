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
        "//app/interface/main/captcha/cmd:all-srcs",
        "//app/interface/main/captcha/conf:all-srcs",
        "//app/interface/main/captcha/dao:all-srcs",
        "//app/interface/main/captcha/http:all-srcs",
        "//app/interface/main/captcha/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
