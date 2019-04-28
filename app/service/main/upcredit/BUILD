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
        "//app/service/main/upcredit/cmd:all-srcs",
        "//app/service/main/upcredit/common/election:all-srcs",
        "//app/service/main/upcredit/common/fsm:all-srcs",
        "//app/service/main/upcredit/conf:all-srcs",
        "//app/service/main/upcredit/dao/account:all-srcs",
        "//app/service/main/upcredit/dao/monitor:all-srcs",
        "//app/service/main/upcredit/dao/upcrmdao:all-srcs",
        "//app/service/main/upcredit/http:all-srcs",
        "//app/service/main/upcredit/mathutil:all-srcs",
        "//app/service/main/upcredit/model:all-srcs",
        "//app/service/main/upcredit/rpc/client:all-srcs",
        "//app/service/main/upcredit/rpc/server:all-srcs",
        "//app/service/main/upcredit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
