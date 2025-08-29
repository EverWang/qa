import { test, expect } from '@playwright/test';

test.describe('管理端功能测试', () => {
  test.beforeEach(async ({ page }) => {
    // 访问管理端登录页面
    await page.goto('http://localhost:3001');
  });

  test('管理员登录功能', async ({ page }) => {
    // 检查登录页面元素
    await expect(page.locator('input[placeholder*="用户名"]')).toBeVisible();
    await expect(page.locator('input[placeholder*="密码"]')).toBeVisible();
    await expect(page.locator('button:has-text("登录")')).toBeVisible();

    // 输入登录信息
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    
    // 点击登录按钮
    await page.click('button:has-text("登录")');
    
    // 等待登录成功，检查是否跳转到仪表板
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    // 验证登录成功后的页面元素
    await expect(page.locator('text=仪表板')).toBeVisible({ timeout: 5000 });
  });

  test('题目管理页面访问', async ({ page }) => {
    // 先登录
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    // 访问题目管理页面
    await page.click('text=题目管理');
    await page.waitForURL('**/questions', { timeout: 5000 });
    
    // 检查题目管理页面元素
    await expect(page.locator('text=题目管理')).toBeVisible();
    await expect(page.locator('button:has-text("添加题目")')).toBeVisible();
  });

  test('新增题目功能', async ({ page }) => {
    // 先登录并进入题目管理页面
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    await page.click('text=题目管理');
    await page.waitForURL('**/questions', { timeout: 5000 });
    
    // 点击添加题目按钮
    await page.click('button:has-text("添加题目")');
    
    // 等待跳转到题目创建页面
    await page.waitForURL('**/questions/create', { timeout: 5000 });
    await expect(page.locator('text=添加题目')).toBeVisible({ timeout: 5000 });
    
    // 填写题目内容
    await page.fill('textarea[placeholder*="请输入题目内容"]', '这是一个测试题目：Vue.js是什么？');
    
    // 选择题目类型（单选题默认已选中）
    await page.click('input[value="single"]');
    
    // 选择分类
    await page.click('.el-select');
    await page.waitForTimeout(500);
    await page.click('.el-option:first-child');
    
    // 选择难度（简单默认已选中）
    await page.click('input[value="easy"]');
    
    // 填写选项
    await page.fill('input[placeholder="请输入选项内容"]:nth-of-type(1)', '前端框架');
    await page.fill('input[placeholder="请输入选项内容"]:nth-of-type(2)', '后端框架');
    
    // 设置正确答案
    await page.click('input[value="A"]');
    
    // 填写解析
    await page.fill('textarea[placeholder*="请输入题目解析"]', 'Vue.js是一个用于构建用户界面的渐进式前端框架。');
    
    // 提交表单
    await page.click('button:has-text("保存题目")');
    
    // 验证题目创建成功
    await expect(page.locator('text=题目创建成功')).toBeVisible({ timeout: 5000 });
  });

  test('题目列表显示和乱码检查', async ({ page }) => {
    // 先登录并进入题目管理页面
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    await page.click('text=题目管理');
    await page.waitForURL('**/questions', { timeout: 5000 });
    
    // 等待题目列表加载
    await page.waitForSelector('table', { timeout: 10000 });
    
    // 检查是否有题目数据
    const rows = await page.locator('table tbody tr').count();
    console.log(`题目列表行数: ${rows}`);
    
    if (rows > 0) {
      // 检查第一行数据是否有乱码
      const firstRowText = await page.locator('table tbody tr:first-child').textContent();
      console.log(`第一行内容: ${firstRowText}`);
      
      // 检查是否包含乱码字符（如问号、方块等）
      expect(firstRowText).not.toMatch(/[\uFFFD\u25A1\u25A0]/); // 检查替换字符和方块字符
      
      // 检查中文字符是否正常显示
      if (firstRowText && /[\u4e00-\u9fa5]/.test(firstRowText)) {
        console.log('检测到中文字符，验证显示正常');
      }
    }
  });

  test('题目编辑功能', async ({ page }) => {
    // 先登录并进入题目管理页面
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    await page.click('text=题目管理');
    await page.waitForURL('**/questions', { timeout: 5000 });
    
    // 等待题目列表加载
    await page.waitForSelector('table', { timeout: 10000 });
    
    // 检查是否有题目可以编辑
    const editButtons = await page.locator('button:has-text("编辑")').count();
    
    if (editButtons > 0) {
      // 点击第一个编辑按钮
      await page.click('button:has-text("编辑"):first');
      
      // 等待编辑对话框出现
      await expect(page.locator('text=编辑题目')).toBeVisible({ timeout: 5000 });
      
      // 修改题目标题
      await page.fill('input[placeholder*="题目标题"]', '修改后的题目标题');
      
      // 保存修改
      await page.click('button:has-text("确定")');
      
      // 验证修改成功
      await expect(page.locator('text=修改成功')).toBeVisible({ timeout: 5000 });
    } else {
      console.log('没有找到可编辑的题目');
    }
  });

  test('题目删除功能', async ({ page }) => {
    // 先登录并进入题目管理页面
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    await page.click('text=题目管理');
    await page.waitForURL('**/questions', { timeout: 5000 });
    
    // 等待题目列表加载
    await page.waitForSelector('table', { timeout: 10000 });
    
    // 记录删除前的题目数量
    const initialCount = await page.locator('table tbody tr').count();
    console.log(`删除前题目数量: ${initialCount}`);
    
    if (initialCount > 0) {
      // 点击第一个删除按钮
      await page.click('button:has-text("删除"):first');
      
      // 确认删除
      await page.click('button:has-text("确定")');
      
      // 验证删除成功
      await expect(page.locator('text=删除成功')).toBeVisible({ timeout: 5000 });
      
      // 验证题目数量减少
      await page.waitForTimeout(1000); // 等待列表刷新
      const finalCount = await page.locator('table tbody tr').count();
      console.log(`删除后题目数量: ${finalCount}`);
      
      expect(finalCount).toBeLessThan(initialCount);
    } else {
      console.log('没有找到可删除的题目');
    }
  });
});