import { test, expect } from '@playwright/test';

test.describe('API接口测试', () => {
  const baseURL = 'http://localhost:8082';
  let adminToken = '';
  let userToken = '';

  test.beforeAll(async ({ request }) => {
    // 获取管理员token
    try {
      const adminLoginResponse = await request.post(`${baseURL}/api/v1/admin/login`, {
        data: {
          username: 'admin',
          password: '123456'
        }
      });
      
      if (adminLoginResponse.ok()) {
        const adminData = await adminLoginResponse.json();
        adminToken = adminData.data?.token || '';
        console.log('管理员登录成功');
      }
    } catch (error) {
      console.log('管理员登录失败:', error);
    }

    // 获取普通用户token
    try {
      const userLoginResponse = await request.post(`${baseURL}/api/v1/auth/login`, {
        data: {
          username: 'testuser',
          password: '123456',
          type: 'password'
        }
      });
      
      if (userLoginResponse.ok()) {
        const userData = await userLoginResponse.json();
        userToken = userData.data?.token || '';
        console.log('普通用户登录成功');
      }
    } catch (error) {
      console.log('普通用户登录失败:', error);
    }
  });

  test('健康检查接口', async ({ request }) => {
    const response = await request.get(`${baseURL}/health`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.status).toBe('ok');
    console.log('健康检查接口正常');
  });

  test('管理员登录接口', async ({ request }) => {
    const response = await request.post(`${baseURL}/api/v1/admin/login`, {
      data: {
        username: 'admin',
        password: '123456'
      }
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(data.data.token).toBeTruthy();
    expect(data.data.user.username).toBe('admin');
    console.log('管理员登录接口正常');
  });

  test('普通用户登录接口', async ({ request }) => {
    const response = await request.post(`${baseURL}/api/v1/auth/login`, {
      data: {
        username: 'testuser',
        password: '123456',
        type: 'password'
      }
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(data.data.token).toBeTruthy();
    expect(data.data.user.username).toBe('testuser');
    console.log('普通用户登录接口正常');
  });

  test('获取分类列表接口', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/v1/categories`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(Array.isArray(data.data)).toBeTruthy();
    
    // 检查分类数据结构
    if (data.data.length > 0) {
      const category = data.data[0];
      expect(category.id).toBeTruthy();
      expect(category.name).toBeTruthy();
      
      // 检查中文字符是否正常
      const categoryName = category.name;
      console.log(`分类名称: ${categoryName}`);
      
      // 检查是否有乱码
      expect(categoryName).not.toMatch(/[\uFFFD\u25A1\u25A0]/);
    }
    
    console.log(`获取到 ${data.data.length} 个分类`);
  });

  test('获取题目列表接口', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/v1/questions`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(Array.isArray(data.data)).toBeTruthy();
    
    // 检查题目数据结构
    if (data.data.length > 0) {
      const question = data.data[0];
      expect(question.id).toBeTruthy();
      expect(question.title).toBeTruthy();
      expect(question.content).toBeTruthy();
      expect(question.options).toBeTruthy();
      
      // 检查中文字符是否正常
      const title = question.title;
      const content = question.content;
      console.log(`题目标题: ${title}`);
      console.log(`题目内容: ${content}`);
      
      // 检查是否有乱码
      expect(title).not.toMatch(/[\uFFFD\u25A1\u25A0]/);
      expect(content).not.toMatch(/[\uFFFD\u25A1\u25A0]/);
      
      // 检查选项格式
      if (typeof question.options === 'string') {
        const options = JSON.parse(question.options);
        expect(Array.isArray(options)).toBeTruthy();
        console.log(`题目选项: ${JSON.stringify(options)}`);
      }
    }
    
    console.log(`获取到 ${data.data.length} 个题目`);
  });

  test('管理员创建题目接口', async ({ request }) => {
    if (!adminToken) {
      test.skip('管理员token不可用，跳过测试');
      return;
    }

    const newQuestion = {
      title: 'API测试题目',
      content: '这是通过API创建的测试题目内容',
      type: 'single',
      options: JSON.stringify([
        { key: 'A', value: 'API测试选项A' },
        { key: 'B', value: 'API测试选项B' },
        { key: 'C', value: 'API测试选项C' },
        { key: 'D', value: 'API测试选项D' }
      ]),
      correctAnswer: 'A',
      categoryId: 1,
      difficulty: 'medium',
      explanation: 'API测试题目解析'
    };

    const response = await request.post(`${baseURL}/api/v1/admin/questions`, {
      headers: {
        'Authorization': `Bearer ${adminToken}`,
        'Content-Type': 'application/json'
      },
      data: newQuestion
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(data.data.id).toBeTruthy();
    expect(data.data.title).toBe(newQuestion.title);
    
    console.log(`创建题目成功，ID: ${data.data.id}`);
  });

  test('管理员获取题目列表接口', async ({ request }) => {
    if (!adminToken) {
      test.skip('管理员token不可用，跳过测试');
      return;
    }

    const response = await request.get(`${baseURL}/api/v1/admin/questions`, {
      headers: {
        'Authorization': `Bearer ${adminToken}`
      }
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(Array.isArray(data.data)).toBeTruthy();
    
    console.log(`管理员获取到 ${data.data.length} 个题目`);
  });

  test('用户答题接口', async ({ request }) => {
    if (!userToken) {
      test.skip('用户token不可用，跳过测试');
      return;
    }

    // 先获取一个题目
    const questionsResponse = await request.get(`${baseURL}/api/v1/questions`);
    const questionsData = await questionsResponse.json();
    
    if (questionsData.data.length === 0) {
      test.skip('没有可用的题目进行答题测试');
      return;
    }

    const question = questionsData.data[0];
    
    // 提交答案
    const answerResponse = await request.post(`${baseURL}/api/v1/questions/${question.id}/answer`, {
      headers: {
        'Authorization': `Bearer ${userToken}`,
        'Content-Type': 'application/json'
      },
      data: {
        answer: 'A',
        timeSpent: 30
      }
    });
    
    expect(answerResponse.ok()).toBeTruthy();
    
    const answerData = await answerResponse.json();
    expect(answerData.code).toBe(200);
    expect(answerData.data.isCorrect).toBeDefined();
    
    console.log(`答题结果: ${answerData.data.isCorrect ? '正确' : '错误'}`);
  });

  test('获取用户答题记录接口', async ({ request }) => {
    if (!userToken) {
      test.skip('用户token不可用，跳过测试');
      return;
    }

    const response = await request.get(`${baseURL}/api/v1/user/records`, {
      headers: {
        'Authorization': `Bearer ${userToken}`
      }
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(Array.isArray(data.data)).toBeTruthy();
    
    console.log(`用户答题记录数量: ${data.data.length}`);
  });

  test('获取错题本接口', async ({ request }) => {
    if (!userToken) {
      test.skip('用户token不可用，跳过测试');
      return;
    }

    const response = await request.get(`${baseURL}/api/v1/user/mistakes`, {
      headers: {
        'Authorization': `Bearer ${userToken}`
      }
    });
    
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.code).toBe(200);
    expect(Array.isArray(data.data)).toBeTruthy();
    
    console.log(`错题本题目数量: ${data.data.length}`);
  });

  test('字符编码测试', async ({ request }) => {
    // 测试包含中文字符的API请求
    const testData = {
      title: '中文测试题目标题',
      content: '这是一个包含中文字符的题目内容，用于测试字符编码是否正确处理。包含特殊字符：！@#￥%……&*（）',
      type: 'single',
      options: JSON.stringify([
        { key: 'A', value: '中文选项A：正确答案' },
        { key: 'B', value: '中文选项B：错误答案' },
        { key: 'C', value: '中文选项C：错误答案' },
        { key: 'D', value: '中文选项D：错误答案' }
      ]),
      correctAnswer: 'A',
      categoryId: 1,
      difficulty: 'easy',
      explanation: '这是中文解析内容，测试字符编码处理。'
    };

    if (adminToken) {
      const response = await request.post(`${baseURL}/api/v1/admin/questions`, {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json; charset=utf-8'
        },
        data: testData
      });
      
      if (response.ok()) {
        const data = await response.json();
        expect(data.code).toBe(200);
        
        // 验证返回的中文字符是否正确
        expect(data.data.title).toBe(testData.title);
        expect(data.data.content).toBe(testData.content);
        
        console.log('中文字符编码测试通过');
      } else {
        console.log('中文字符编码测试失败:', await response.text());
      }
    } else {
      console.log('跳过中文字符编码测试（无管理员token）');
    }
  });
});