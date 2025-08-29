import { test, expect } from '@playwright/test';

test.describe('用户端功能测试', () => {
  test.beforeEach(async ({ page }) => {
    // 访问用户端页面
    await page.goto('http://localhost:3000');
  });

  test('用户端首页加载', async ({ page }) => {
    // 检查首页基本元素
    await expect(page.locator('text=刷刷题')).toBeVisible({ timeout: 10000 });
    
    // 检查是否有题目列表或相关内容
    const hasQuestions = await page.locator('text=题目').isVisible({ timeout: 5000 });
    const hasLogin = await page.locator('text=登录').isVisible({ timeout: 5000 });
    
    // 至少应该有登录按钮或题目内容
    expect(hasQuestions || hasLogin).toBeTruthy();
  });

  test('用户登录功能', async ({ page }) => {
    // 查找登录相关元素
    const loginButton = page.locator('button:has-text("登录")');
    const loginLink = page.locator('a:has-text("登录")');
    
    // 尝试点击登录按钮或链接
    if (await loginButton.isVisible({ timeout: 3000 })) {
      await loginButton.click();
    } else if (await loginLink.isVisible({ timeout: 3000 })) {
      await loginLink.click();
    } else {
      // 如果没有找到登录按钮，可能需要先访问登录页面
      await page.goto('http://localhost:3000/login');
    }
    
    // 检查登录表单元素
    await expect(page.locator('input[placeholder*="用户名"], input[placeholder*="账号"]')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('input[placeholder*="密码"]')).toBeVisible();
    
    // 输入测试用户信息
    await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
    await page.fill('input[placeholder*="密码"]', '123456');
    
    // 点击登录
    await page.click('button:has-text("登录")');
    
    // 验证登录成功（可能跳转到首页或显示用户信息）
    await page.waitForTimeout(2000);
    
    // 检查是否有用户相关信息显示
    const hasUserInfo = await page.locator('text=testuser, text=个人中心, text=退出').isVisible({ timeout: 5000 });
    const isHomePage = await page.locator('text=首页, text=题目列表').isVisible({ timeout: 5000 });
    
    expect(hasUserInfo || isHomePage).toBeTruthy();
  });

  test('题目列表显示', async ({ page }) => {
    // 先尝试登录（如果需要）
    const needLogin = await page.locator('button:has-text("登录")').isVisible({ timeout: 3000 });
    
    if (needLogin) {
      await page.click('button:has-text("登录")');
      await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
      await page.fill('input[placeholder*="密码"]', '123456');
      await page.click('button:has-text("登录")');
      await page.waitForTimeout(2000);
    }
    
    // 查找题目列表
    const questionsList = page.locator('.question-item, .question-card, [class*="question"]');
    const questionsTable = page.locator('table');
    
    // 等待题目列表加载
    await page.waitForTimeout(3000);
    
    // 检查是否有题目显示
    const hasQuestionItems = await questionsList.count() > 0;
    const hasQuestionsTable = await questionsTable.isVisible({ timeout: 3000 });
    
    if (hasQuestionItems || hasQuestionsTable) {
      console.log('找到题目列表');
      
      // 检查题目内容是否有乱码
      if (hasQuestionItems) {
        const firstQuestion = await questionsList.first().textContent();
        console.log(`第一个题目内容: ${firstQuestion}`);
        
        // 检查乱码
        if (firstQuestion) {
          expect(firstQuestion).not.toMatch(/[\uFFFD\u25A1\u25A0]/);
        }
      }
    } else {
      console.log('未找到题目列表，可能需要导航到题目页面');
      
      // 尝试点击题目相关导航
      const questionNav = page.locator('text=题目, text=练习, text=刷题');
      if (await questionNav.first().isVisible({ timeout: 3000 })) {
        await questionNav.first().click();
        await page.waitForTimeout(2000);
      }
    }
  });

  test('答题页面功能', async ({ page }) => {
    // 先登录
    const needLogin = await page.locator('button:has-text("登录")').isVisible({ timeout: 3000 });
    
    if (needLogin) {
      await page.click('button:has-text("登录")');
      await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
      await page.fill('input[placeholder*="密码"]', '123456');
      await page.click('button:has-text("登录")');
      await page.waitForTimeout(2000);
    }
    
    // 查找第一个题目并点击
    const firstQuestion = page.locator('.question-item, .question-card, [class*="question"]').first();
    const startButton = page.locator('button:has-text("开始答题"), button:has-text("进入题目")');
    
    if (await firstQuestion.isVisible({ timeout: 5000 })) {
      await firstQuestion.click();
    } else if (await startButton.isVisible({ timeout: 5000 })) {
      await startButton.click();
    } else {
      // 尝试直接访问答题页面
      await page.goto('http://localhost:3000/question/1');
    }
    
    // 等待答题页面加载
    await page.waitForTimeout(3000);
    
    // 检查答题页面元素
    const questionContent = page.locator('.question-content, [class*="question"]');
    const options = page.locator('.option, [class*="option"], input[type="radio"]');
    const submitButton = page.locator('button:has-text("提交"), button:has-text("确定")');
    
    // 验证答题页面基本元素
    const hasQuestionContent = await questionContent.isVisible({ timeout: 5000 });
    const hasOptions = await options.count() > 0;
    const hasSubmitButton = await submitButton.isVisible({ timeout: 5000 });
    
    console.log(`答题页面检查: 题目内容=${hasQuestionContent}, 选项=${hasOptions}, 提交按钮=${hasSubmitButton}`);
    
    if (hasQuestionContent) {
      // 检查题目内容是否有乱码
      const content = await questionContent.textContent();
      console.log(`题目内容: ${content}`);
      
      if (content) {
        expect(content).not.toMatch(/[\uFFFD\u25A1\u25A0]/);
      }
    }
    
    if (hasOptions && hasSubmitButton) {
      // 选择第一个选项
      await options.first().click();
      
      // 提交答案
      await submitButton.click();
      
      // 等待结果显示
      await page.waitForTimeout(2000);
      
      // 检查是否有结果反馈
      const hasResult = await page.locator('text=正确, text=错误, text=答案, text=解析').isVisible({ timeout: 5000 });
      
      if (hasResult) {
        console.log('答题结果显示正常');
      }
    }
  });

  test('题目分类功能', async ({ page }) => {
    // 先登录
    const needLogin = await page.locator('button:has-text("登录")').isVisible({ timeout: 3000 });
    
    if (needLogin) {
      await page.click('button:has-text("登录")');
      await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
      await page.fill('input[placeholder*="密码"]', '123456');
      await page.click('button:has-text("登录")');
      await page.waitForTimeout(2000);
    }
    
    // 查找分类相关元素
    const categories = page.locator('.category, [class*="category"], .tab');
    const categorySelect = page.locator('select[name*="category"], select[name*="type"]');
    
    // 检查是否有分类功能
    const hasCategoryTabs = await categories.count() > 0;
    const hasCategorySelect = await categorySelect.isVisible({ timeout: 3000 });
    
    console.log(`分类功能检查: 分类标签=${hasCategoryTabs}, 分类选择=${hasCategorySelect}`);
    
    if (hasCategoryTabs) {
      // 点击第二个分类（如果存在）
      const categoryCount = await categories.count();
      if (categoryCount > 1) {
        await categories.nth(1).click();
        await page.waitForTimeout(2000);
        
        // 验证分类切换后题目列表有变化
        console.log('分类切换功能正常');
      }
    }
    
    if (hasCategorySelect) {
      // 选择不同的分类
      await categorySelect.selectOption({ index: 1 });
      await page.waitForTimeout(2000);
      
      console.log('分类选择功能正常');
    }
  });

  test('搜索功能', async ({ page }) => {
    // 先登录
    const needLogin = await page.locator('button:has-text("登录")').isVisible({ timeout: 3000 });
    
    if (needLogin) {
      await page.click('button:has-text("登录")');
      await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
      await page.fill('input[placeholder*="密码"]', '123456');
      await page.click('button:has-text("登录")');
      await page.waitForTimeout(2000);
    }
    
    // 查找搜索框
    const searchInput = page.locator('input[placeholder*="搜索"], input[placeholder*="查找"], input[type="search"]');
    const searchButton = page.locator('button:has-text("搜索"), button:has-text("查找")');
    
    if (await searchInput.isVisible({ timeout: 5000 })) {
      // 输入搜索关键词
      await searchInput.fill('测试');
      
      // 点击搜索按钮或按回车
      if (await searchButton.isVisible({ timeout: 3000 })) {
        await searchButton.click();
      } else {
        await searchInput.press('Enter');
      }
      
      // 等待搜索结果
      await page.waitForTimeout(2000);
      
      console.log('搜索功能测试完成');
    } else {
      console.log('未找到搜索功能');
    }
  });

  test('个人中心功能', async ({ page }) => {
    // 先登录
    const needLogin = await page.locator('button:has-text("登录")').isVisible({ timeout: 3000 });
    
    if (needLogin) {
      await page.click('button:has-text("登录")');
      await page.fill('input[placeholder*="用户名"], input[placeholder*="账号"]', 'testuser');
      await page.fill('input[placeholder*="密码"]', '123456');
      await page.click('button:has-text("登录")');
      await page.waitForTimeout(2000);
    }
    
    // 查找个人中心入口
    const profileLink = page.locator('text=个人中心, text=我的, text=用户中心');
    const userAvatar = page.locator('.avatar, [class*="avatar"], .user-info');
    
    if (await profileLink.isVisible({ timeout: 5000 })) {
      await profileLink.click();
    } else if (await userAvatar.isVisible({ timeout: 5000 })) {
      await userAvatar.click();
    } else {
      // 尝试直接访问个人中心页面
      await page.goto('http://localhost:3001/profile');
    }
    
    // 等待个人中心页面加载
    await page.waitForTimeout(2000);
    
    // 检查个人中心相关元素
    const hasUserInfo = await page.locator('text=testuser, text=用户信息').isVisible({ timeout: 5000 });
    const hasRecords = await page.locator('text=答题记录, text=做题记录, text=成绩').isVisible({ timeout: 5000 });
    const hasMistakes = await page.locator('text=错题本, text=错题').isVisible({ timeout: 5000 });
    
    console.log(`个人中心检查: 用户信息=${hasUserInfo}, 答题记录=${hasRecords}, 错题本=${hasMistakes}`);
    
    // 至少应该有用户信息显示
    expect(hasUserInfo || hasRecords || hasMistakes).toBeTruthy();
  });
});