Page({

  /**
   * 页面的初始数据
   */
  data: {
    showPop: false,  // 新增待办弹出框
    taskName: "",   // 任务名称
    priority: '0',  // 任务优先级
    todos: [],  

    inputValue: '',
    startX: 0,
  },

  showPopup() {
    this.setData({ showPop: true });
  },

  onClose() {
    this.setData({ showPop: false });
  },

  onTaskNameChange(event) {
    this.setData({ taskName: event.detail });
  },

  // 优先级改变
  onPriorityChange(event) {
    this.setData({
      priority: event.detail,
    });
  },

  // 输入框内容变化
  inputChange(e) {
    this.setData({
      inputValue: e.detail.value
    })
  },

  // 拉取某天数据
  getTaskList() {
    wx.showLoading({
      title: '加载中...'
    })

    this.data.todos = []

    wx.request({
      url: 'http://localhost:18888/task/list',
      method: 'get',
      data: {
        "user_id": 1, 
        "date": '2024-12-22',
        "type": 0,
        "year": '2024',
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer 112233mrkleo' // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          wx.showToast({
            title: '获取数据失败',
            icon: 'none'
          })
          return 
        }

        let tasks = []
        for(let item of res.data) {
          tasks.push({
            id: item.task_id,
            content: item.name, 
            showDelete: false,
          })
        }

        this.setData({
          todos: tasks,
        });
      },
      fail: (err) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      },
      complete: () => {
        wx.hideLoading()
      }
    })
  },

  // 添加待办事项
  addTodo() {
    if (!this.data.taskName.trim()) {
      wx.showToast({
        title: '请输入内容',
        icon: 'none'
      })
      return
    }

    // 请求接口添加任务
    wx.showLoading({
      title: '加载中...'
    })

    let priority = parseInt(this.data.priority)

    wx.request({
      url: 'http://localhost:18888/task/create',
      method: 'POST',
      data: {
        "user_id": 1, 
        "name": this.data.taskName.trim(),
        "type": 0,
        "priority": priority,
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer 112233mrkleo' // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          wx.showToast({
            title: '获取数据失败',
            icon: 'none'
          })
          return 
        }

        this.data.showPop = false
        this.data.taskName = ""
        this.data.priority = "0"
        // 重新拉取数据
        this.getTaskList()        
      },
      fail: (err) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      },
      complete: () => {
        wx.hideLoading()
      }
    })
  },

  // 完成某项任务
  onFinishTask() {
    wx.showLoading({
      title: '加载中...'
    })
    // wx.request({
    //   url: 'http://localhost:18888/task/finished',
    //   method: 'post',
    //   data: {
    //     "task_id": 1, 
    //   },
    //   header: {
    //     'content-type': 'application/json',
    //     'Authorization': 'Bearer 112233mrkleo' // 如果需要token
    //   },
    //   success: (res) => {
    //      this.getTaskList()
    //   },
    //   fail: (err) => {
    //     wx.showToast({
    //       title: '网络错误',
    //       icon: 'none'
    //     })
    //   },
    //   complete: () => {
    //     wx.hideLoading()
    //   }
    // })
  },

  // 格式化时间
  formatTime(date) {
    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const day = date.getDate()
    const hour = date.getHours()
    const minute = date.getMinutes()
    return `${year}-${month}-${day} ${hour}:${minute}`
  },

  // 触摸开始
  touchStart(e) {
    const touch = e.touches[0]
    this.setData({
      startX: touch.clientX
    })
  },

  // 触摸移动
  touchMove(e) {
    const touch = e.touches[0]
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    const item = todos[index]
    
    // 计算移动距离
    let moveLength = touch.clientX - this.data.startX
    const deleteWidth = 320

    if (moveLength < 0) { // 左滑
      item.showDelete = true
    } else { // 右滑
      item.showDelete = false
    }

    this.setData({
      todos
    })
  },

  // 触摸结束
  touchEnd(e) {
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    const item = todos[index]

    if (item.showDelete) {
      item.showDelete = true
    } else {
      item.showDelete = false
    }

    this.setData({
      todos
    })
  },

  // 点击隐藏删除按钮
  hideDelete(e) {
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    todos[index].showDelete = false
    this.setData({
      todos
    })
  },

  // 删除待办事项
  deleteTodo(e) {
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    todos.splice(index, 1)
    this.setData({
      todos
    })
    wx.showToast({
      title: '删除成功',
      icon: 'success'
    })
  },

  // 修改时间
  editTime(e) {
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    
    wx.showModal({
      title: '提示',
      content: '确定要更新时间为当前时间吗？',
      success: (res) => {
        if (res.confirm) {
          todos[index].time = this.formatTime(new Date())
          todos[index].showDelete = false
          this.setData({
            todos
          })
          wx.showToast({
            title: '更新成功',
            icon: 'success'
          })
        }
      }
    })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.getTaskList()
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})