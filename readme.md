# wallpaper-filter

## API
### 获取指定目录的图片列表
```
GET
/image_list
```
Query string:
```
dir=${dir_path} // 根目录下相对路径，前端硬编码给几个入口
hwoperator=${hw_operator}   // 宽高比运算符，>=: ge; <=: le
hwratio=${hw_ratio} // 宽高比，string, 后端parse成float32
pagesize=${page_size}   // 分页加载，单页数量
pagenum=${page_num} // 分页加载，页码offset
pagenum=${page_num} // 分页加载，页码offset
```

### 下载图片到本地目录
```
POST
/image
```
Param:
```
file=${file_full_path}  // 从根目录开始计算的相对路径
targetdir=${target_dir} // 存放的目录，从根目录开始计算的相对路径；可选，不填为后端默认下载路径
```