# 顶层设计

感谢 [gin-vue-admin](https://www.gin-vue-admin.com/)

---

## **已部署上线**

 [在线预览](http://iadmin.xyz/)

测试账号：12312312312 测试密码：123

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

## 已完成功能 （当前已更新版本1.3.0）

- [x] 多点登录拦截
- [x] 分页表格（列表页，详情页，新增页）（基础的增删改查）（完成时间：2021年6月15日）
- [x] 带快捷选项的日期范围选择器（完成时间：2021年6月15日）
- [x] 优化滚动条样式（完成时间：2021年6月15日）
- [x] 设置自动刷新（1s，3s，5s)（完成时间：2021年6月15日）
- [x] 权限管理（完成时间：2021年6月17日)
- [x] 不同权限不同菜单（完成时间：2021年6月17日）
- [x] 动态路由（完成时间：2021年6月18日）
- [x] 操作日志记录（完成时间：2021年6月16日）
- [x] 动态菜单（完成时间：2021年6月18日）
- [x] 错误页面（404）
- [x] 菜单优化、图标可自定义、可自定义菜单顺序、不同角色可以设置不同首页（完成时间：2021年6月30日）
- [x] 刷新保存菜单选择、展开情况（完成时间：2021年7月1日）
- [x] 优化选择标签页，会对应显示其对应的菜单，和展开其父菜单（完成时间：2021年7月1日）
- [x] 解决刷新返回根目录的问题（完成时间：2021年7月1日）
- [x] 路由面包屑（完成时间：2021年7月2日）
- [x] Casbin 权限管理（完成时间：2021年7月2日）
- [x] 富文本编辑器（前端）（完成时间：2021年7月3日）
- [x] 路由切换动画（完成时间：2021年7月4日）
- [x] 上传文件（基础）（完成时间：2021年7月5日）
- [x] 批量导出
- [x] 批量导入
- [x] 个人中心
- [x] 前端密码加密

