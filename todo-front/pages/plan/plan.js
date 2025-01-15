const recorderManager = wx.getRecorderManager();

let api = require('../../common/api');

Page({

  /**
   * 页面的初始数据
   */
  data: {
    messages: [
      { id: Date.now(), content: 'Hello,我是您的AI创建计划助手,您可以输入‘英语四级学习计划’、‘云南大理旅行’来创建计划~', sender: 'ai' }
    ], // 消息数组
    inputValue: '', // 输入框内容
    aiAvatar: '../../images/ai_suggest.png', // AI 头像（需替换为真实路径）
    userAvatar: '../../images/default-avatar.png', // 用户头像（需替换为真实路径）
    conversationID: "", // 会话ID, 首次创建不传递

    recording: false,    // 是否正在录音
    cancelRecord: false  // 是否取消录音
  },
  onInput(e) {
    this.setData({ inputValue: e.detail.value });
  },
  sendMessage() {
    const { inputValue, messages } = this.data;
    if (!inputValue.trim()) return;

    // 添加用户消息
    messages.push({ id: Date.now(), content: inputValue, sender: 'user' });
    this.setData({ 
      messages, 
      inputValue: '' 
    });

    // 模拟调用 AI 接口
    this.getAIResponse(inputValue);
  },
  getAIResponse(userMessage) {
    const { messages } = this.data;
    const that = this;
  
    // 显示“AI 正在输入...”的提示
    messages.push({ id: Date.now(), content: '正在输入...', sender: 'ai', isTyping: true });
    this.setData({ messages });
  
    // 发起请求
    wx.request({
      url: api.ApiHost + "/im_plan/chat", 
      method: 'POST',
      data: { 
        user_id: api.UserID, 
        query: userMessage,
        conversation_id: this.data.conversationID,
      },
      header: {
        Accept: 'text/event-stream', // 指定 SSE 数据流格式
      },
      success(res) {
        if (!res || !res.data) {
          return 
        }

        const sseData = res.data; // 接收流式数据
        let currentContent = '';
        const messages = that.data.messages;

        // 解析 SSE 返回的数据流
        sseData.split('\n\n').forEach((line) => {
          if (line.startsWith('data:')) {
            try {
              // 提取并解析 JSON 数据
              const jsonString = line.replace('data:', '').trim();
              const jsonData = JSON.parse(jsonString);

              // 判断是否为 `agent_message` 事件并提取 answer 字段
              if (jsonData.event === 'agent_message' && jsonData.answer) {
                currentContent += jsonData.answer;

                // 更新当前 "正在输入..." 的消息
                const updatedMessages = messages.map((msg) =>
                  msg.isTyping
                    ? { ...msg, content: currentContent } // 更新 AI 正在输入的内容
                    : msg
                );
                that.setData({ messages: updatedMessages });
              }

              // 解析到conversation_id, 就记录下
              if (jsonData.event === 'agent_message' && jsonData.conversation_id) {
                that.setData({ conversationID: jsonData.conversation_id });
              }

            } catch (error) {
              console.error('解析 SSE 数据失败:', error);
            }
          }
        });

        that.setData({
          messages: that.data.messages.map((msg) => {
            if (msg.isTyping) {
              if (msg.content.includes("```")) {
                // 提取计划内容(通过```xxx```包裹)
                let contArr = msg.content.split("```")
                let planCont = ""
                if (contArr.length > 2) {
                  planCont = contArr[1]
                } else {
                  return { id: Date.now(), content: currentContent, sender: 'ai' } 
                }
                return { id: Date.now(), content: currentContent, sender: 'ai', genPlan: true, planCont: planCont } 
              } 
              return { id: Date.now(), content: currentContent, sender: 'ai' } 
            } else {
              return msg
            }
          }),
        });
      },
      fail() {
        // 请求失败处理
        const failedMessages = that.data.messages.filter((msg) => !msg.isTyping);
        failedMessages.push({ id: Date.now(), content: 'AI 回复失败，请重试。', sender: 'ai' });
        that.setData({ messages: failedMessages });
      },
    });

  },

  // 应用Ai计划
  onConfirmApply(e) {
    wx.request({
      url: api.ApiHost + '/im_plan/apply',
      method: 'post',
      data: {
        "user_id": api.UserID, 
        ai_gen_cont: e.detail
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (res.data == 'ok') {
          wx.showToast({
            title: '创建成功',
          })
        }
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

  // 加载历史消息
  async loadHistoryMessages(conversation_id) {
    let converDetail = await this.getConversationDetail(conversation_id);

    let {messages} = this.data;
    for(let item of converDetail.data) {
      messages.push({ id: Date.now(), content: item.query, sender: 'user' });
      messages.push({ id: Date.now(), content: item.answer, sender: 'ai' });
    }

    this.setData({
      messages: this.data.messages.map((msg) => {
        if (msg.content.includes("```")) {
          // 提取计划内容(通过```xxx```包裹)
          let contArr = msg.content.split("```")
          let planCont = ""
          if (contArr.length > 2) {
            planCont = contArr[1]
          } else {
            return { id: msg.id, content: msg.content, sender: msg.sender } 
          }
          return { id: msg.id, content: msg.content, sender: msg.sender, genPlan: true, planCont: planCont } 
        } 
        return { id: msg.id, content: msg.content, sender: msg.sender  } 
      }),
    });
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

  // 开始录音
  startRecord(e) {
    // 先获取录音权限
    wx.authorize({
      scope: 'scope.record',
      success: () => {
        this.setData({
          recording: true,
          cancelRecord: false
        });
        
        // 开始录音的配置
        const options = {
          duration: 60000,
          sampleRate: 16000,
          numberOfChannels: 1,
          encodeBitRate: 96000,
          format: 'mp3',
          frameSize: 50
        }
        
        recorderManager.start(options);
      },
      fail: () => {
        wx.showModal({
          title: '提示',
          content: '需要您的录音权限',
          success: (res) => {
            if (res.confirm) {
              wx.openSetting();
            }
          }
        });
      }
    });
  },

  // 结束录音
  endRecord() {
    if (!this.data.recording) return;
    
    this.setData({
      recording: false
    });
    
    if (this.data.cancelRecord) {
      recorderManager.stop();
      return;
    }
    
    recorderManager.stop();
  },

  // 将语音转为文字
  translateVoice(tempFilePath) {
    wx.showLoading({
      title: '正在识别...'
    });

    // 使用微信原生的语音识别接口
    wx.uploadFile({
      url: 'YOUR_SERVER_URL', // 你的服务器接口地址
      filePath: tempFilePath,
      name: 'file',
      success: (res) => {
        const result = JSON.parse(res.data);
        if (result.success) {
          this.setData({
            inputValue: result.text
          });
        } else {
          wx.showToast({
            title: '识别失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        console.error('语音识别失败:', err);
        wx.showToast({
          title: '识别失败',
          icon: 'none'
        });
      },
      complete: () => {
        wx.hideLoading();
      }
    });
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    // 监听录音结束事件
    recorderManager.onStop((res) => {
      if (this.data.cancelRecord) {
        console.log('取消录音');
        return;
      }
      
      this.translateVoice(res.tempFilePath);
    });

    // 获取会话id参数(从其他页面跳转过来)
    if (options.id) {
      this.setData({ conversationID: options.id });
      this.loadHistoryMessages(options.id);
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