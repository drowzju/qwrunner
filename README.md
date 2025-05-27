 go build  



 # 视觉理解

 .\qwrunner.exe -i https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg -m qwen-vl-max -q “图片里有几个生物？”


 # 视觉推理

 .\qwrunner.exe -m qvq-max -i "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg" -q "请解答这道题" -s


```
PS D:\code\qwrunner> .\qwrunner.exe -m qvq-max -i "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg" -q "请解答这道题" -s
**答案：**

1. **长方体**
   - **表面积**：\( 2 \times (4 \times 3 + 4 \times 2 + 3 \times 2) = 2 \times 26 = 52 \, \text{cm}^2 \)
   - **体积**：\( 4 \times 3 \times 2 = 24 \, \text{cm}^3 \)

2. **正方体**
   - **表面积**：\( 6 \times 3^2 = 6 \times 9 = 54 \, \text{cm}^2 \)
   - **体积**：\( 3^3 = 27 \, \text{cm}^3 \)

**总结：**
- 长方体的表面积为 \( 52 \, \text{cm}^2 \)，体积为 \( 24 \, \text{cm}^3 \)。
- 正方体的表面积为 \( 54 \, \text{cm}^2 \)，体积为 \( 27 \, \text{cm}^3 \)。
PS D:\code\qwrunner>

```

