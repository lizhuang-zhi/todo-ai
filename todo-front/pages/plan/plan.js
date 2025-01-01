// pages/plan.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    messages: [], // 存储消息
    inputValue: '', // 用户输入的内容
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
    if (!this.data.showConfirm) {
      // 阻止页面卸载，显示弹框
      this.setData({ showConfirm: true });
      wx.showModal({
        title: '确认退出',
        content: '您确定要离开此页面吗？',
        success: (res) => {
          if (res.confirm) {
            // 用户确认退出，继续返回
            wx.navigateBack();
          } else {
            // 用户取消退出，什么都不做
            this.setData({ showConfirm: false });
          }
        },
      });
    }
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