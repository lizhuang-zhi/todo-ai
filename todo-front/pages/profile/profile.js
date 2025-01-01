let api = require('../../common/api');

Page({
  data: {
  },

  // 跳转数据页面
  onTapMore() {
    wx.navigateTo({
      url: '/pages/data/data',
    })
  },

  async onLoad() {
    this.getUserInfo();
  },

  getUserInfo() {
    // 获取用户信息逻辑
    wx.getUserProfile({
      desc: '用于完善会员资料',
      success: (res) => {
        this.setData({
          userInfo: res.userInfo
        });
      }
    });
  },

  // 页面相关事件处理函数
  onPullDownRefresh() {
  },
    /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
    if (typeof this.getTabBar === 'function' &&
      this.getTabBar()) {
      this.getTabBar().setData({
        selected: 1
      })
    }
  },
});