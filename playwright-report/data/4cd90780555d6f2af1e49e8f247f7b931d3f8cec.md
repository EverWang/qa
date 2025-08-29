# Page snapshot

```yaml
- generic [ref=e4]:
  - generic [ref=e5]:
    - img [ref=e7]
    - heading "刷刷题" [level=1] [ref=e9]
    - paragraph [ref=e10]: 专业的在线刷题平台
  - generic [ref=e11]:
    - generic [ref=e12]:
      - generic [ref=e13]: 用户名
      - textbox "用户名" [ref=e14]: testuser
    - generic [ref=e15]:
      - generic [ref=e16]: 密码
      - textbox "密码" [ref=e17]: "123456"
    - button "登录中..." [disabled] [ref=e18]:
      - img [ref=e19]
      - generic [ref=e21]: 登录中...
  - paragraph [ref=e23]:
    - text: 登录即表示同意
    - link "用户协议" [ref=e24] [cursor=pointer]:
      - /url: "#"
    - text: 和
    - link "隐私政策" [ref=e25] [cursor=pointer]:
      - /url: "#"
```