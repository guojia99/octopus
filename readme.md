# octopus
> - 这是一个拥有各种各样基于基础库开发的一个工具库
> - 本仓库基于go1.18开发
> - 基本上每个工具都是带锁的，保证并发安全

# 依赖
```
> - json 依赖 `github.com/json-iterator/go`
```

# 项目列表

### 1、OMap
- [`omap`](./omap): 基于`list.List` 封装的一个map, 支持按key排序，支持拓展的map相关函数
- `Map` 为any原型
- `OMap`为泛型
- 使用方法参见 [`map example`](./_example/map)
### 2、OArray
- [`oarray`](./oarray):  带锁的数组列表，包含一些实用方法
- `Array` 为any 原型
- `OArray` 为泛型
- 使用方法参见 [`array example`](./_example/array) 
### 3、OFunc
- [`ofunc`](./ofunc): go函数操作集合
  - Try-Catch: 错误捕捉操作, 使用方法见[`try_catch_example`](./_example/func/try_catch.go)
