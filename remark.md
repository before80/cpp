stdlib.md比较特殊，
会找不到：ldiv、mblen、mbtowc、wctomb、mbstowcs、wcstombs，原因是div所在页面，会出现 多个h2标签导致找不到,
故只能手动！

stdio.md比较特殊，会提示找不到很多函数，原因fopen所在页面，出现了“文件访问标记”这个h3标签（已处理，程序自动替换了）

stdbit.md比较特殊，都没有中文页面

stdarg.md很奇怪，也找不到 va_end 和  va_start，待查找什么问题！<- 出现了“展开值”这个h3标签（已处理，程序自动替换了）

math.md很奇怪，会提示找不到很多函数，待查找什么问题！