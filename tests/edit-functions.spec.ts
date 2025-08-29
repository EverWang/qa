import { test, expect } from '@playwright/test'

test.describe('编辑功能测试', () => {
  test.beforeEach(async ({ page }) => {
    // 登录管理员
    await page.goto('http://localhost:3000/login')
    await page.fill('input[placeholder*="用户名"]', 'admin')
    await page.fill('input[placeholder*="密码"]', '123456')
    await page.click('button:has-text("登录")')
    await page.waitForURL('**/dashboard', { timeout: 10000 })
  })

  test('分类编辑功能 - 设置上级分类', async ({ page }) => {
    // 进入分类管理页面
    await page.click('text=分类管理')
    await page.waitForURL('**/categories', { timeout: 5000 })
    
    // 等待分类列表加载
    await page.waitForSelector('table', { timeout: 10000 })
    
    // 查找一个可编辑的分类（非大类）
    const editButtons = await page.locator('button:has-text("编辑")')
    const count = await editButtons.count()
    
    if (count > 0) {
      // 点击第一个编辑按钮
      await editButtons.first().click()
      
      // 等待编辑对话框出现
      await expect(page.locator('text=编辑分类')).toBeVisible({ timeout: 5000 })
      
      // 选择上级分类
      await page.click('.el-select')
      await page.waitForSelector('.el-option', { timeout: 3000 })
      
      // 选择第一个大类作为上级分类
      const options = await page.locator('.el-option')
      if (await options.count() > 1) {
        await options.nth(1).click() // 跳过"无（作为大类）"选项
      }
      
      // 保存修改
      await page.click('button:has-text("确定")')
      
      // 验证保存成功
      await expect(page.locator('text=更新成功')).toBeVisible({ timeout: 5000 })
      
      console.log('分类编辑功能测试通过')
    } else {
      console.log('没有找到可编辑的分类')
    }
  })

  test('分类编辑功能 - 设置为大类', async ({ page }) => {
    // 进入分类管理页面
    await page.click('text=分类管理')
    await page.waitForURL('**/categories', { timeout: 5000 })
    
    // 等待分类列表加载
    await page.waitForSelector('table', { timeout: 10000 })
    
    // 查找一个可编辑的分类
    const editButtons = await page.locator('button:has-text("编辑")')
    const count = await editButtons.count()
    
    if (count > 0) {
      // 点击第一个编辑按钮
      await editButtons.first().click()
      
      // 等待编辑对话框出现
      await expect(page.locator('text=编辑分类')).toBeVisible({ timeout: 5000 })
      
      // 选择"无（作为大类）"
      await page.click('.el-select')
      await page.waitForSelector('.el-option', { timeout: 3000 })
      await page.click('text=无（作为大类）')
      
      // 保存修改
      await page.click('button:has-text("确定")')
      
      // 验证保存成功
      await expect(page.locator('text=更新成功')).toBeVisible({ timeout: 5000 })
      
      console.log('分类设置为大类功能测试通过')
    } else {
      console.log('没有找到可编辑的分类')
    }
  })

  test('题目编辑功能', async ({ page }) => {
    // 进入题目管理页面
    await page.click('text=题目管理')
    await page.waitForTimeout(1000) // 等待菜单展开
    await page.click('text=题目列表')
    await page.waitForURL('**/questions', { timeout: 5000 })
    
    // 等待题目列表加载
    await page.waitForSelector('table', { timeout: 10000 })
    
    // 查找一个可编辑的题目
    const editButtons = await page.locator('button:has-text("编辑")')
    const count = await editButtons.count()
    
    if (count > 0) {
      // 点击第一个编辑按钮
      await editButtons.first().click()
      
      // 等待编辑页面加载
      await expect(page.locator('text=编辑题目')).toBeVisible({ timeout: 5000 })
      
      // 修改题目内容
      const originalContent = await page.inputValue('textarea[placeholder*="题目内容"]')
      const newContent = originalContent + ' (已修改)'
      await page.fill('textarea[placeholder*="题目内容"]', newContent)
      
      // 保存修改
      await page.click('button:has-text("保存题目")')
      
      // 验证保存成功
      await expect(page.locator('text=题目更新成功')).toBeVisible({ timeout: 5000 })
      
      console.log('题目编辑功能测试通过')
    } else {
      console.log('没有找到可编辑的题目')
    }
  })
})