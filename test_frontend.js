// 前端功能测试脚本
const axios = require('axios');

// 测试配置
const API_BASE = 'http://localhost:8080';
const ADMIN_WEB = 'http://localhost:5173';
const MINIPROGRAM = 'http://localhost:5174';

/**
 * 测试后端API连接
 */
async function testBackendAPI() {
    console.log('\n=== 测试后端API连接 ===');
    try {
        const response = await axios.get(`${API_BASE}/health`);
        console.log('✅ 后端API连接成功:', response.status);
        return true;
    } catch (error) {
        console.log('❌ 后端API连接失败:', error.message);
        return false;
    }
}

/**
 * 测试管理员登录API
 */
async function testAdminLogin() {
    console.log('\n=== 测试管理员登录 ===');
    try {
        const response = await axios.post(`${API_BASE}/api/v1/admin/login`, {
            username: 'admin',
            password: '123456'
        });
        console.log('✅ 管理员登录成功');
        return response.data.data?.token || response.data.token;
    } catch (error) {
        console.log('❌ 管理员登录失败:', error.response?.data || error.message);
        return null;
    }
}

/**
 * 测试获取分类列表
 */
async function testGetCategories(token) {
    console.log('\n=== 测试获取分类列表 ===');
    try {
        const response = await axios.get(`${API_BASE}/api/v1/admin/categories`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        console.log('✅ 获取分类列表成功，分类数量:', response.data.data?.length || 0);
        return response.data.data;
    } catch (error) {
        console.log('❌ 获取分类列表失败:', error.response?.data || error.message);
        return [];
    }
}

/**
 * 测试获取题目列表
 */
async function testGetQuestions(token) {
    console.log('\n=== 测试获取题目列表 ===');
    try {
        const response = await axios.get(`${API_BASE}/api/v1/admin/questions`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        console.log('✅ 获取题目列表成功，题目数量:', response.data.total || 0);
        return response.data.data;
    } catch (error) {
        console.log('❌ 获取题目列表失败:', error.response?.data || error.message);
        return [];
    }
}

/**
 * 测试前端页面可访问性
 */
async function testFrontendPages() {
    console.log('\n=== 测试前端页面可访问性 ===');
    
    // 测试管理端
    try {
        const adminResponse = await axios.get(ADMIN_WEB);
        console.log('✅ 管理端页面可访问:', ADMIN_WEB);
    } catch (error) {
        console.log('❌ 管理端页面不可访问:', error.message);
    }
    
    // 测试小程序端
    try {
        const miniprogramResponse = await axios.get(MINIPROGRAM);
        console.log('✅ 小程序端页面可访问:', MINIPROGRAM);
    } catch (error) {
        console.log('❌ 小程序端页面不可访问:', error.message);
    }
}

/**
 * 主测试函数
 */
async function runTests() {
    console.log('🚀 开始前端功能自测...');
    
    // 测试后端API
    const backendOk = await testBackendAPI();
    if (!backendOk) {
        console.log('\n❌ 后端服务未启动，请先启动后端服务');
        return;
    }
    
    // 测试管理员登录
    const token = await testAdminLogin();
    if (!token) {
        console.log('\n❌ 管理员登录失败，无法继续测试');
        return;
    }
    
    // 测试API功能
    await testGetCategories(token);
    await testGetQuestions(token);
    
    // 测试前端页面
    await testFrontendPages();
    
    console.log('\n🎉 前端功能自测完成！');
    console.log('\n📋 测试结果总结:');
    console.log('- 后端API: ✅ 正常');
    console.log('- 管理员登录: ✅ 正常');
    console.log('- 管理端页面: http://localhost:5173');
    console.log('- 小程序端页面: http://localhost:5174');
    console.log('\n💡 建议手动在浏览器中访问以上页面进行进一步测试');
}

// 运行测试
runTests().catch(console.error);