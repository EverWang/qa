# TODO:

- [x] fix-admin-category-disabled-filter: 修复管理端分类列表搜索条件选择禁用时筛选结果不正确的问题 - 后端筛选逻辑正常，前端传参正确 (priority: High)
- [x] fix-admin-question-add-save-error: 修复添加题目保存时报错参数错误的问题（编辑题目保存正常）- 已添加缺失的type字段 (priority: High)
- [x] fix-admin-question-type-filter: 修复题目管理页面题目类型筛选条件不起作用的问题 - 已添加type参数到API请求 (priority: High)
- [x] fix-admin-operation-logs-username: 修复操作日志中操作人全部显示为Unknown，应该显示正确用户名的问题 - 已修复LogOperation函数根据用户角色查询正确的用户名 (priority: High)
- [x] fix-mini-homepage-categories-display: 修复小程序首页没有显示分类信息，需要显示分类并可点击进入该分类查看题目 - 已修复分类筛选逻辑 (priority: High)
- [x] fix-mini-category-page-dropdown: 修复分类页面顶部分类下拉框没有数据可选，缺少筛选条件的问题 - 已修复下拉框显示所有分类并实现正确的筛选逻辑 (priority: High)
- [x] fix-mini-guest-login-navigation: 修复点击游客体验弹出登录成功但没有跳转到正确页面（如错题本、个人中心）的问题 - 已修复路由守卫允许游客用户访问需要认证的页面 (priority: High)
