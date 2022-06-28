# octopus
> - 这是一个拥有各种各样基于基础库开发的一个工具库
> - 本仓库基于go1.18开发
> - 本仓库不依赖其他第三方库

# 项目列表

### 1、OMap

- [`omap`](./omap): 基于`list.List` 封装的一个map, 支持按key排序，支持拓展的map相关函数，支持锁,

- `Map` 为any原型

  - ```go
    package main
    
    func main(){
     	m := NewMap()
        value := []byte{1, 2, 3}
        m.Set("key", value)
        data, err := m.Get("key")
        fmt.Println("data -> %v, err -> %v", data, err)
    }
    ```

  - 

- `OMap`为泛型