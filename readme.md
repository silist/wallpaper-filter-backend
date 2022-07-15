# wallpaper-filter

## API
### 获取指定目录的图片列表
```
GET
/image_list
```
Query String:
```
dir=${dir_path} // 根目录下相对路径，前端硬编码给几个入口
hwoperator=${hw_operator}   // 宽高比运算符，>=: gte; <=: lte
hwratio=${hw_ratio} // 宽高比，string, 后端parse成float32
pagesize=${page_size}   // 分页加载，单页数量
pagenum=${page_num} // 分页加载，页码offset
```

### 发送图片
```
GET
/image
```
Query String:
```
path=${rel_path} // 图片在根目录下相对路径
```

### 下载图片到本地目录
```
POST
/image
```
Param:
```
path=${rel_path}  // 从根目录开始计算的相对路径
```