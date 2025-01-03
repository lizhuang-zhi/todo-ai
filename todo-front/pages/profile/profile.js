let api = require('../../common/api');
let utils = require('../../utils/utils');

Page({
  data: {
    conversations: [],  // 会话数据
  },

  // 跳转数据页面
  onTapMore() {
    wx.navigateTo({
      url: '/pages/data/data',
    })
  },

  async onLoad() {
    this.getUserInfo();
    this.getInitData();  // 获取会话数据
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

  // 获取会话数据
  async getInitData() {
    // TODO: 不处理分页，暂时仅加载20条数据
    // 同步获取会话列表
    let cvstResp = await this.getConversations();
    let conversations = cvstResp.data;

    let result = [];

    // 同步获取会话详情
    for (let i = 0; i < conversations.length; i++) {
      let csn = conversations[i];
      let csnDetail = await this.getConversationDetail(csn.id);

      let {lastMsg, canShare} = this.getLastMessage(csnDetail.data)

      result.push({
        id: csn.id,
        name: csn.name,
        lastMsg: lastMsg,  // 聊天对话中机器人返回的最后一条消息
        canShare: canShare,  // 是否可以分享
      });
    }

    this.setData({
      conversations: result,
    })
  },

  // 同步获取会话列表
  async getConversations() {
    return new Promise((resolve, reject) => {
      wx.request({
        url: api.ApiHost + '/im_plan/conversations',
        method: 'get',
        data: {
          pinned: true,
          last_id: "",
          limit: 20,   // TODO:暂时写死，不分页
        },
        header: {
          'content-type': 'application/json',
          'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
        },
        success: (res) => {
          if (!res || !res.data) {
            return reject('获取会话列表失败')
          }
          resolve(res.data)
        },
        fail: (err) => {
          reject('获取会话列表失败')
        },
      })
    })
  },

  // 同步获取会话详情
  async getConversationDetail(conversationId) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: api.ApiHost + '/im_plan/messages/',
        method: 'get',
        data: {
          conversation_id: conversationId,
          first_id: "",
          limit: 20,   // TODO:暂时写死，不分页
        },
        header: {
          'content-type': 'application/json',
          'Authorization': 'Bearer'+ api.AuthKey // 如果需要token
        },
        success: (res) => {
          if (!res ||!res.data) {
            return reject('获取会话详情失败')
          }
          resolve(res.data)
        },
        fail: (err) => {
          reject('获取会话详情失败')
        },
      })
    })
  },

  // 获取会话详情最后一条消息
  getLastMessage(messages) {
    if (!messages || messages.length === 0) {
      return {
        lastMsg: '',
        canShare: false,
      }
    }

    // 获取最后一条消息
    let lastConversation = messages[messages.length - 1];
    let lastMsg = utils.getShortStr(lastConversation, 30);
    
    let canShare = false;
    for(let item of messages) {
      if (item.answer != "" && item.answer.includes("```")) {
        canShare = true;
        break;
      }
    }

    return {
      lastMsg: lastMsg,
      canShare: canShare,
    }
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