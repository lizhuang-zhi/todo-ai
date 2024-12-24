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
      pagePath: "/pages/index/index",
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
  // 添加页面显示的生命周期函数
  pageLifetimes: {
    show: function() {
      // 获取当前页面路径
      const pages = getCurrentPages();
      const currentPage = pages[pages.length - 1];
      const route = '/' + currentPage.route;
      
      // 根据当前页面路径设置选中状态
      const selected = this.data.list.findIndex(item => item.pagePath === route);
      if (selected !== -1) {
        this.setData({
          selected
        });
      }
    }
  },
  methods: {
    // tabbar装换
    switchTab(e) {
      const data = e.currentTarget.dataset;
      const url = data.path;
      
      // 先设置选中状态
      this.setData({
        selected: data.index
      });
      
      // 再进行页面跳转
      wx.switchTab({
        url,
        fail: (err) => {
          console.error('切换失败：', err);
        }
      });
    },
    // 跳转发布页面
    toRelease() {
      wx.navigateTo({
        url: '/pages/logs/logs',
      })
    }
  }
})