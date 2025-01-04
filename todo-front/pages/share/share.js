// pages/share/share.js
Page({

  data: {
    userInfo: {
      avatarUrl: '/images/default-avatar.png',
      nickName: '李白'
    },
    hasUserInfo: false,

    buttonText: "",  // 按钮文案(老用户: 分享计划, 新用户: 一键应用)
  },

  onGetVip() {
    // 如果还没有用户信息，先获取用户信息
    if (!this.data.hasUserInfo) {
      this.getUserProfile()
      return
    }

    wx.showModal({
      title: '领取成功',
      content: '恭喜获得7天高级会员体验',
      showCancel: false,
      success: (res) => {
        if (res.confirm) {
          // 跳转到首页或其他页面
          wx.switchTab({
            url: '/pages/index/index'
          })
        }
      }
    })
  },

  // 添加获取用户信息的方法
  getUserProfile() {
    wx.getUserProfile({
      desc: '用于完善会员资料',
      success: (res) => {
        console.log(res)
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    })
  },

  onShareAppMessage() {
    return {
      title: 'xxxxxx, 一起提升效率吧',
      path: '/pages/share/share',
      imageUrl: '/images/default-avatar.png'
    }
  },

  // 一键应用模板
  onJoinPlan() {
    if (this.data.buttonText === "分享计划") {
      // 调用微信的分享功能
      this.onShareAppMessage()
      return
    } else if (this.data.buttonText === "一键应用") {
     this.onApplyPlan()
     return 
    }
  },

  // 应用计划(新用户)
  onApplyPlan() {
    wx.showModal({
      title: '应用成功',
      content: '已成功加入日程~',
      showCancel: false
    })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    // TODO: 用于测试展示
    // options.id = undefined

    if (options.id) {
      this.setData({ buttonText: "分享计划" })  // 老用户
    } else {
      this.setData({ buttonText: "一键应用" })  // 新用户
    } 
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