// pages/plan.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    messages: [], // 消息数组
    inputValue: '', // 输入框内容
    aiAvatar: '../../images/ai_suggest.png', // AI 头像（需替换为真实路径）
    userAvatar: '../../images/default-avatar.png', // 用户头像（需替换为真实路径）
  },
  onInput(e) {
    this.setData({ inputValue: e.detail.value });
  },
  sendMessage() {
    const { inputValue, messages } = this.data;
    if (!inputValue.trim()) return;

    // 添加用户消息
    messages.push({ id: Date.now(), content: inputValue, sender: 'user' });
    this.setData({ messages, inputValue: '' });

    // 模拟调用 AI 接口
    this.getAIResponse(inputValue);
  },
  getAIResponse(userMessage) {
    // 模拟 AI 回复
    wx.request({
      url: 'https://your-dify-api-endpoint', // 替换为你的 Dify API 地址
      method: 'POST',
      data: { message: userMessage },
      success: (res) => {
        const aiReply = res.data.reply || 'AI 未能理解您的问题';
        const { messages } = this.data;
        messages.push({ id: Date.now(), content: aiReply, sender: 'ai' });
        this.setData({ messages });
      },
      fail: () => {
        const { messages } = this.data;
        messages.push({ id: Date.now(), content: 'AI 回复失败，请重试。', sender: 'ai' });
        this.setData({ messages });
      },
    });
  },

  // 点击返回
  onBackClick() {
    wx.showModal({
      title: '确认退出?',
      content: '历史对话可在个人页面查看',
      success: (res) => {
        if (res.confirm) {
          // 确认退出
          wx.navigateBack();
        }
      },
    });
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {

  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {
  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  }
})