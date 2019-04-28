package(default_visibility = ["//visibility:public"])

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
        "//app/service/main/assist/model/assist:all-srcs",
        "//app/service/main/assist/model/message:all-srcs",
        "//app/service/main/assist/model/notify:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
