# Page snapshot

```yaml
- generic [ref=e4]:
  - generic [ref=e5]:
    - img [ref=e7]
    - heading "刷刷题" [level=1] [ref=e10]
    - paragraph [ref=e11]: 专业的在线刷题平台
  - generic [ref=e12]:
    - generic [ref=e13]:
      - generic [ref=e14]: 用户名
      - textbox "用户名" [ref=e15]: testuser
    - generic [ref=e16]:
      - generic [ref=e17]: 密码
      - textbox "密码" [ref=e18]: "123456"
    - button "登录中..." [disabled] [ref=e19]:
      - img [ref=e20]
      - generic [ref=e22]: 登录中...
  - paragraph [ref=e24]:
    - text: 登录即表示同意
    - link "用户协议" [ref=e25] [cursor=pointer]:
      - /url: "#"
    - text: 和
    - link "隐私政策" [ref=e26] [cursor=pointer]:
      - /url: "#"
```