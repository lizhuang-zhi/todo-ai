let api = require('../../common/api');
let utils = require('../../utils/utils');

Page({

  /**
   * 页面的初始数据
   */
  data: {
    selectDate: "",  // 当前选中日期
    holidays: {}, // 用于存储从后端获取的节假日数据

    showPop: false,  // 新增待办弹出框
    taskName: "",   // 任务名称
    priority: '0',  // 任务优先级
    editID: "",   // 编辑时的任务id
    todos: [],  
    startX: 0,

    unFinishPng: "../../images/unfinish.png",
    finishedPng: "../../images/finished.png",
    showAiSuggestPop: false,    // 是否展示AI建议框
    aiSuggestCont: false,   // AI建议内容
  },

  // 选择日期(日历组件)
  onSelectDate(e) {
    let detailDate = e.detail
    this.setData({
      selectDate: utils.formatDate(detailDate)
    })

    this.getTaskList()
  },

  showPopup() {
    this.setData({ showPop: true });
  },

  onClose() {
    this.setData({ 
      showPop: false,
      editID: "", 
     });
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

  // 拉取某天数据
  getTaskList() {
    wx.showLoading({
      title: '加载中...'
    })

    this.setData({
      todos: []
    })

    wx.request({
      url: api.ApiHost + '/task/list',
      method: 'get',
      data: {
        "user_id": 1, 
        "date": this.data.selectDate,
        "type": 0,
        // TODO: 改为年获取
        "year": '2024',
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          return 
        }

        let tasks = []
        for(let item of res.data) {
          if (item.progress == 1) {
            continue
          }

          tasks.push({
            id: item.task_id,
            content: item.name, 
            priority: item.priority,
            todoPng: this.data.unFinishPng,  // 统一为未完成
            showDelete: false,
            aiSuggestion: item.ai_suggestion,  // ai建议内容
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
      url: api.ApiHost + '/task/create',
      method: 'POST',
      data: {
        "user_id": 1, 
        "name": this.data.taskName.trim(),
        "type": 0,
        "priority": priority,
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          wx.showToast({
            title: '获取数据失败',
            icon: 'none'
          })
          return 
        }

        this.setData({
          taskName: '', 
          priority: '0',
          showPop: false,
        })
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
  onFinishTask(e) {
    // 先将图片标记为完成, 等200ms后删除
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    todos[index].todoPng = this.data.finishedPng
    this.setData({
      todos
    })

    let timer = setTimeout(() => {
      todos.splice(index, 1)
      this.setData({
        todos
      })
      clearTimeout(timer);
    }, 100)

    const id = e.currentTarget.dataset.id
    wx.request({
      url: api.ApiHost + '/task/finished',
      method: 'post',
      data: {
        "task_id": id, 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
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

  // 删除待办事项
  deleteTodo(e) {
    const index = e.currentTarget.dataset.index
    const todos = this.data.todos
    todos.splice(index, 1)
    this.setData({
      todos
    })

    const id = e.currentTarget.dataset.id
    wx.request({
      url: api.ApiHost + '/task/delete',
      method: 'post',
      data: {
        "task_id": id, 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
          this.getTaskList()
      },
      fail: (err) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      },
      complete: () => {
      }
    })
  },

  // 打开修改弹出框
  editPop(e) {
    let item = e.currentTarget.dataset.item
    this.setData({
      taskName: item.content, 
      priority: item.priority.toString(),
      editID: item.id,
      showPop: true,
    })
  }, 

  // 修改待办
  editTodo() {
    wx.showLoading({
      title: '修改中...'
    })
    wx.request({
      url: api.ApiHost + '/task/update',
      method: 'post',
      data: {
        "task_id": this.data.editID, 
        "name": this.data.taskName,
        "priority": parseInt(this.data.priority),
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          wx.showToast({
            title: '获取数据失败',
            icon: 'none'
          })
          return 
        }

        // 重置数据
        this.setData({
          taskName: '', 
          priority: '0',
          editID: '',
          showPop: false,
        })
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

  // 展示AI建议
  showAISuggestion(e) {
    let item = e.currentTarget.dataset.item
    this.setData({
      aiSuggestCont: item.aiSuggestion,
      showAiSuggestPop: true,  
    })
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

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.getHolidays()  // 获取日历数据
    this.initDate()   // 初始化日期
    this.getTaskList()
  },

  // vant日历组件所需函数
  formatter(day) {
    // 将日期转换为 YYYY-MM-DD 格式
    const year = day.date.getFullYear();
    const month = String(day.date.getMonth() + 1).padStart(2, '0');
    const date = String(day.date.getDate()).padStart(2, '0');
    const dateStr = `${year}-${month}-${date}`;

    // 检查是否是节假日
    const holiday = this.data.holidays[dateStr];
    if (holiday) {
      // 法定节假日显示红色标记
      if (holiday.isHoliday) {
        day.topInfo = holiday.name;
        day.bottomInfo = '休';
        day.className = 'holiday';
        day.text = day.date.getDate();
      } else {
        // 普通节日只显示名称
        day.topInfo = holiday.name;
      }
    }

    // 处理日期范围选择的起始和结束标记
    if (day.type === 'start') {
      day.bottomInfo = '开始';
    } else if (day.type === 'end') {
      day.bottomInfo = '结束';
    }

    // 标记今天
    const today = new Date();
    if (today.getFullYear() === year && 
        today.getMonth() === day.date.getMonth() && 
        today.getDate() === day.date.getDate()) {
      day.text = '今天';
    }

    return day;
  },

  // 获取节假日数据
  async getHolidays() {
    let date = new Date()

    wx.request({
      url: api.ApiHost + '/calendar/data',
      method: 'get',
      data: {
        "year": date.getFullYear(), 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey 
      },
      success: (res) => {
        if (!res || !res.data) {
          return 
        }

        // 转换数据格式为 { 'YYYY-MM-DD': holidayObject }
        const holidays = {};
        res.data.forEach(item => {
          holidays[item.date] = item;
        });

        this.setData({ 
          holidays,
          formatter: this.formatter
        });
      },
      fail: (err) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      },
      complete: () => {
      }
    })
  },

  initDate() {  
    // 获取当天日期
    let nowDate = new Date()
    this.setData({
      selectDate: utils.formatDate(nowDate)
    })
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
    if (typeof this.getTabBar === 'function' &&
      this.getTabBar()) {
      this.getTabBar().setData({
        selected: 0
      })
    }
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