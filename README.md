# 顶层设计

感谢 [gin-vue-admin](https://www.gin-vue-admin.com/)

---

## **已部署上线**

 [在线预览](http://iadmin.xyz/)

## 前端框架

1. Vue3.0，官方文档：[https://v3.cn.vuejs.org/guide/migration/introduction.html#概览](https://v3.cn.vuejs.org/guide/migration/introduction.html#%E6%A6%82%E8%A7%88)
2. UI样式库，Ant Design，官网文档：[https://2x.antdv.com/components/overview-cn/](https://2x.antdv.com/components/overview-cn/)
3. NProgress，github地址：[https://github.com/rstacruz/nprogress](https://github.com/rstacruz/nprogress)
4. vue-persist，官网文档：[https://www.npmjs.com/package/vuex-persist](https://www.npmjs.com/package/vuex-persist)

---

## 后端框架

> Gin 是一个用 Go (Golang) 编写的 HTTP web 框架。 它是一个类似于 martini 但拥有更好性能的 API 框架, 优于 httprouter，速度提高了近 40 倍。如果你需要极好的性能，使用 Gin 吧。官方文档：[https://gin-gonic.com/zh-cn/docs/](https://gin-gonic.com/zh-cn/docs/)

- gorm，中文文档：[https://learnku.com/docs/gorm/v2/index/9728](https://learnku.com/docs/gorm/v2/index/9728)

    注意gorm有v1和v2版本，v2版本迁移到[gorm.io/gorm](http://gorm.io/gorm)

- Viper，github地址：[https://github.com/spf13/viper](https://github.com/spf13/viper)

    viper 是一个配置解决方案

- Redis

    目前只有开启多点登录拦截时才会启动Redis

- JWT
- Swagger2

    接口说明文档   swag init

- Zap记录日志

    下载安装：go get [go.uber.org/zap](http://go.uber.org/zap)

    学习网址：[https://www.liwenzhou.com/posts/Go/use_zap_in_gin/](https://www.liwenzhou.com/posts/Go/use_zap_in_gin/)

    file-rotatelogs:[github.com/lestrrat-go/file-rotatelogs](http://github.com/lestrrat-go/file-rotatelogs)

---

## 前端功能（页面）

1. 登录页
2. 注册页
3. 框架
    1. logo
    2. 顶部用户
    3. 左侧菜单栏
    4. Tab列表
    5. 主体内容框架
        1. 头部
        2. 主体内容
        3. 尾部
4. cookie存储（弃用，使用本地存储local）

以上完成时间🕝2021年6月13日

---

## 已完成功能

- [x]  多点登录拦截
- [x]  分页表格（列表页，详情页，新增页）（基础的增删改查）（完成时间：2021年6月15日）
- [x]  带快捷选项的日期范围选择器（完成时间：2021年6月15日）
- [x]  优化滚动条样式（完成时间：2021年6月15日）
- [x]  设置自动刷新（1s，3s，5s)（完成时间：2021年6月15日）
- [x]  权限管理（完成时间：2021年6月17日)
- [x]  不同权限不同菜单（完成时间：2021年6月17日）
- [x]  动态路由（Q：刷新页面动态路由失效？）（完成时间：2021年6月18日）
- [x]  操作日志记录（完成时间：2021年6月16日）
- [x]  接口请求记录（完成时间：2021年6月16日）
- [x]  菜单和路由结合？（完成时间：2021年6月18日）
- [x]  错误页面（404）

---

## 以后计划扩展功能

- [ ]  接口权限配置（Casbin配置）
- [ ]  Q：刷新页面会跳转首页
- [ ]  菜单优化
    - [ ]  不同角色设置不同的首页
    - [ ]  菜单展示顺序可设置
    - [ ]  菜单图标可设置
- [ ]  批量导出
- [ ]  批量导入
- [ ]  动态配置路由
- [ ]  密码明文传值？
- [ ]  Q：菜单收起显示问题
- [ ]  上传文件（上传到哪？本地还是服务器？）（文件病毒问题，上传文件安全）（图片压缩，视频截取封面）
- [ ]  服务器情况
- [ ]  邮件发送
- [ ]  短信发送
- [ ]  前端面包屑路由导航