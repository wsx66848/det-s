# SVG生成

### 构建

```bat
rem windows
set GOOS=linux
go build .
```

```sh
# linux
go build .
```

将 `det-s` `compoent` `vue/dist` 复制到 `10.3.220.200`上的 `/home/byr/fsc/TSD` 文件夹,并启动 `det-s`

### 测试

`go run . -d` 启动服务

客户端上传jpg文件 `curl 127.0.0.1:8003 -F 'file=@test.jpg'` ,服务器会调用 `fsc/detect.py` ，并将输出成svg字符串返回到浏览器，在 `10.3.220.200` 上执行太慢了,为了调试方便,开发时 `detect.py` 简单输出 `result.txt` 里面的内容,我们直接使用这个值即可.

### 流程

从收到一个http请求道返回结果,主要处理流程如下

1. 接收file，交给py处理，获取输出字符串类似 `bboxes: [ 55 323  80 343   0], label_text: netport|1.00`

    参见 `main.go@call()`

1. 将这些信息转换成 `Panel`里面的目标 `Target`

    参见 `panel.go@Panel::Add()`

1. py返回结果中,网口/光口 有`1`口`2`口`4`口,结果不准确(例如`6`口识别成`4+2`,上下同一组的无法识别),而且不方便统一处理,这里都拆分成单个网口/光口,后面再去想办法组合

    参见 `panel.go@Target::vsplite()`

1. 将目标的位置信息按照 `U` 实际尺寸按比例转化, `1U=435mm*44mm`, 这里统一进行了缩放 `1U=1250px*127px`

    参见 `panel.go@Panel::Format()`

1. 将同种类型的目标尝试对齐，例如某N个网口的y坐标十分接近（极差大于`delta=8px`），那么这N个网口的y坐标求平均值avgY,这些网口的y坐标全部设置为avgY,横坐标同理

    参见 `algorithm.go@TargetSlice::adjust()`

1. 将同种类型目标尝试分组,例如某N个网口之间的距离十分接近（上下两个网口x坐标距离小于 `delta` 划分到一个组, 左右两个网口x坐标距离小于 `delta+Nwidth`划分到一个组）

    参见 `algorithm.go@groupTargets()`

1. 对每一个分组(和没分组的,例如 `BackPlane` ),查询对应的SVG文件,根据坐标进行平移插入到SVG中

    参见 `panel.go@Panel::ToSvg()`

### 问题

1. 网口/光口没完全对齐，或者间距异常

    目标检测输出结果位置误差太大，可以适当增大delta（太大会影响分组，把多个组合到了一个组）值，主要还是看识别准确率，而不是放大程序容忍误差

1. 网口/光口与其编号存在一定的偏移

    对应的svg大小不标准，svg网口平均宽度尽可能等于 `Nwidth=41`

    程序中定义的 `Nwidth` 值不准, 设置 `Nwidth=W*1250/435`, `W` 表示网口标准宽度(mm)， `1250` 表示svg面板宽度(px)，`435` 表示1U宽度(mm)

1. svg面板出现 `?`

    看下出现 `?` 位置是什么，缺失对应的svg图片，在component目录下添加对于的svg，程序启动加上 `-d` 可以看到缺失的具体文件是啥，例如输出日志：`获取svg失败: [network] [1m1.svg]`，则补上对应的文件即可

1. 程序只能识别厂商、网口、光口、指示灯、USB、背板。
