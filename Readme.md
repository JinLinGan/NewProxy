## 编译前准备

目前只测试过使用 `mingw64` 进行编译。

从 [github](https://github.com/niXman/mingw-builds-binaries/releases/tag/12.1.0-rt_v10-rev3) 下载 `x86_64-12.1.0-release-win32-seh-rt_v10-rev3.7z`

解压到后配置环境变量 `PATH` 中加入 `C:\xxxxxxx\mingw64\bin`

## 编译 

```golang
go build
```

