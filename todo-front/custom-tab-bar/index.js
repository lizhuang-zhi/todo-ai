Component({
  data: {
    selected: 0,
    color: "#7A7E83",
    selectedColor: "#450003",
    /* 补全list数组 */
    list: [{
      pagePath: "/pages/achievement/achievement",
      iconPath: "/images/achievement.png",
      selectedIconPath: "/images/achievement-select.png",
      text: "成就"
    }, {
      pagePath: "/pages/todo/todo",
      iconPath: "/images/todo.png",
      selectedIconPath: "/images/todo-select.png",
      text: "Todo"
    }, {
      pagePath: "/pages/profile/profile",
      iconPath: "/images/profile.png",
      selectedIconPath: "/images/profile-select.png",
      text: "我的"
    }]
  },
  attached() {
  },
  methods: {
    // tabbar装换
    switchTab(e) {
      const data = e.currentTarget.dataset
      const url = data.path
      wx.switchTab({url})
      this.setData({
        selected: data.index
      })
    },
    // 跳转发布页面
    toRelease() {
      wx.navigateTo({
        url: '/pages/logs/logs',
      })
      console.log('跳转->注意事项页面');
    }
    
  }
})