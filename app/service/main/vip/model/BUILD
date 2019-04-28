load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "associate_banner.go",
        "associate_bind.go",
        "associate_ele_api.go",
        "associate_model.go",
        "associate_old.go",
        "associate_open.go",
        "associate_order.go",
        "associate_panel.go",
        "associate_params.go",
        "bp.go",
        "coupon.go",
        "dialog.go",
        "http.go",
        "jointly.go",
        "mail_api.go",
        "model.go",
        "notify.go",
        "oldvip.go",
        "order.go",
        "panel.go",
        "params.go",
        "pay.go",
        "platform.go",
        "point.go",
        "privilege.go",
        "resource.go",
        "rpc.go",
        "tips.go",
        "user.go",
        "vip.go",
        "welfare.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/service/main/vip/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/vip/model:go_default_library",
        "//app/service/main/coupon/api:go_default_library",
        "//app/service/main/coupon/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "model_proto",
    srcs = ["vip.proto"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/service/main/vip/model",
    proto = ":model_proto",
    tags = ["automanaged"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
    ],
)
