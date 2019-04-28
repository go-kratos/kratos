![Excelize](./excelize.png "Excelize")

# Excelize

[![Build Status](https://travis-ci.org/360EntSecGroup-Skylar/excelize.svg?branch=master)](https://travis-ci.org/360EntSecGroup-Skylar/excelize)
[![Code Coverage](https://codecov.io/gh/360EntSecGroup-Skylar/excelize/branch/master/graph/badge.svg)](https://codecov.io/gh/360EntSecGroup-Skylar/excelize)
[![Go Report Card](https://goreportcard.com/badge/github.com/360EntSecGroup-Skylar/excelize)](https://goreportcard.com/report/github.com/360EntSecGroup-Skylar/excelize)
[![GoDoc](https://godoc.org/github.com/360EntSecGroup-Skylar/excelize?status.svg)](https://godoc.org/github.com/360EntSecGroup-Skylar/excelize)
[![Licenses](https://img.shields.io/badge/license-bsd-orange.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/xuri)

## 简介

Excelize 是 Go 语言编写的用于操作 Office Excel 文档类库，基于 ECMA-376 Office OpenXML 标准。可以使用它来读取、写入由 Microsoft Excel™ 2007 及以上版本创建的 XLSX 文档。相比较其他的开源类库，Excelize 支持写入原本带有图片(表)、透视表和切片器等复杂样式的文档，还支持向 Excel 文档中插入图片与图表，并且在保存后不会丢失文档原有样式，可以应用于各类报表系统中。使用本类库要求使用的 Go 语言为 1.8 或更高版本，完整的 API 使用文档请访问 [godoc.org](https://godoc.org/github.com/360EntSecGroup-Skylar/excelize) 或查看 [参考文档](https://xuri.me/excelize/)。

## 快速上手

### 安装

```go
go get github.com/360EntSecGroup-Skylar/excelize
```

### 创建 Excel 文档

下面是一个创建 Excel 文档的简单例子：

```go
package main

import (
    "fmt"

    "github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
    xlsx := excelize.NewFile()
    // 创建一个工作表
    index := xlsx.NewSheet("Sheet2")
    // 设置单元格的值
    xlsx.SetCellValue("Sheet2", "A2", "Hello world.")
    xlsx.SetCellValue("Sheet1", "B2", 100)
    // 设置工作簿的默认工作表
    xlsx.SetActiveSheet(index)
    // 根据指定路径保存文件
    err := xlsx.SaveAs("./Book1.xlsx")
    if err != nil {
        fmt.Println(err)
    }
}
```

### 读取 Excel 文档

下面是读取 Excel 文档的例子：

```go
package main

import (
    "fmt"

    "github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
    xlsx, err := excelize.OpenFile("./Book1.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    // 获取工作表中指定单元格的值
    cell := xlsx.GetCellValue("Sheet1", "B2")
    fmt.Println(cell)
    // 获取 Sheet1 上所有单元格
    rows := xlsx.GetRows("Sheet1")
    for _, row := range rows {
        for _, colCell := range row {
            fmt.Print(colCell, "\t")
        }
        fmt.Println()
    }
}
```

### 在 Excel 文档中创建图表

使用 Excelize 生成图表十分简单，仅需几行代码。您可以根据工作表中的已有数据构建图表，或向工作表中添加数据并创建图表。

![Excelize](./test/images/chart.png "Excelize")

```go
package main

import (
    "fmt"

    "github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
    categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
    values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
    xlsx := excelize.NewFile()
    for k, v := range categories {
        xlsx.SetCellValue("Sheet1", k, v)
    }
    for k, v := range values {
        xlsx.SetCellValue("Sheet1", k, v)
    }
    xlsx.AddChart("Sheet1", "E1", `{"type":"col3DClustered","series":[{"name":"Sheet1!$A$2","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$2:$D$2"},{"name":"Sheet1!$A$3","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$3:$D$3"},{"name":"Sheet1!$A$4","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$4:$D$4"}],"title":{"name":"Fruit 3D Clustered Column Chart"}}`)
    // 根据指定路径保存文件
    err := xlsx.SaveAs("./Book1.xlsx")
    if err != nil {
        fmt.Println(err)
    }
}

```

### 向 Excel 文档中插入图片

```go
package main

import (
    "fmt"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"

    "github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
    xlsx, err := excelize.OpenFile("./Book1.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    // 插入图片
    err = xlsx.AddPicture("Sheet1", "A2", "./image1.png", "")
    if err != nil {
        fmt.Println(err)
    }
    // 在工作表中插入图片，并设置图片的缩放比例
    err = xlsx.AddPicture("Sheet1", "D2", "./image2.jpg", `{"x_scale": 0.5, "y_scale": 0.5}`)
    if err != nil {
        fmt.Println(err)
    }
    // 在工作表中插入图片，并设置图片的打印属性
    err = xlsx.AddPicture("Sheet1", "H2", "./image3.gif", `{"x_offset": 15, "y_offset": 10, "print_obj": true, "lock_aspect_ratio": false, "locked": false}`)
    if err != nil {
        fmt.Println(err)
    }
    // 保存文件
    err = xlsx.Save()
    if err != nil {
        fmt.Println(err)
    }
}
```

## 社区合作

欢迎您为此项目贡献代码，提出建议或问题、修复 Bug 以及参与讨论对新功能的想法。 XML 符合标准： [part 1 of the 5th edition of the ECMA-376 Standard for Office Open XML](http://www.ecma-international.org/publications/standards/Ecma-376.htm)。

## 致谢

本类库中部分 XML 结构体的定义参考了开源项目：[tealeg/xlsx](https://github.com/tealeg/xlsx).

## 开源许可

本项目遵循 BSD 3-Clause 开源许可协议，访问 [https://opensource.org/licenses/BSD-3-Clause](https://opensource.org/licenses/BSD-3-Clause) 查看许可协议文件。
