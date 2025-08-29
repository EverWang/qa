# Page snapshot

```yaml
- generic [active] [ref=e1]:
  - generic [ref=e4]:
    - generic [ref=e5]:
      - heading "刷刷题管理后台" [level=1] [ref=e6]
      - paragraph [ref=e7]: 请使用管理员账号登录
    - generic [ref=e8]:
      - generic [ref=e12]:
        - img [ref=e15]
        - textbox "请输入用户名" [ref=e17]: admin
      - generic [ref=e21]:
        - img [ref=e24]
        - textbox "请输入密码" [ref=e27]: "123456"
        - img [ref=e30] [cursor=pointer]
      - generic [ref=e35] [cursor=pointer]:
        - generic [ref=e36] [cursor=pointer]:
          - checkbox "记住我"
        - generic [ref=e38] [cursor=pointer]: 记住我
      - button "登录" [ref=e41] [cursor=pointer]:
        - generic [ref=e42] [cursor=pointer]: 登录
    - paragraph [ref=e44]:
      - generic [ref=e45]: 演示账号：admin / 123456
  - alert [ref=e46]:
    - img [ref=e48]
    - paragraph [ref=e50]: 请先登录
  - alert [ref=e51]:
    - img [ref=e53]
    - paragraph [ref=e55]: 登录已过期，请重新登录
  - alert [ref=e56]:
    - img [ref=e58]
    - paragraph [ref=e60]: 登录失败，请检查用户名和密码
  - alert [ref=e61]:
    - img [ref=e63]
    - paragraph [ref=e65]: 登录已过期，请重新登录
  - alert [ref=e66]:
    - img [ref=e68]
    - paragraph [ref=e70]: 登录已过期，请重新登录
```