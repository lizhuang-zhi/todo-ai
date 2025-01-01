Component({
  properties: {
    aiPlanText: {
      type: String,
      observer: function(newVal) {
        if (newVal) {
          this.parseAIContent(newVal)
        }
      }
    }, 
  },

  data: {
    parentTask: {
      name: '',
      date: ''
    },
    sonTasks: []
  },

  methods: {
    parseAIContent(content) {
      // 解析父任务
      const parentMatch = content.match(/\[\[ParentTask\]\]\[add\](.*?)@(.*?)(\n|$)/);
      if(parentMatch) {
        this.setData({
          'parentTask.name': parentMatch[1],
          'parentTask.date': parentMatch[2]
        });
      }

      // 解析子任务
      const sonTaskContent = content.match(/\[\[SonTask\]\](.*?)(\n|$)/);
      if(sonTaskContent) {
        const tasks = sonTaskContent[1].split('|||');
        const sonTasks = [];
        tasks.forEach(task => {
          const match = task.match(/\[add\](.*?)@(.*?)$/);
          if(match) {
            sonTasks.push({
              name: match[1],
              date: match[2]
            });
          }
        });
        this.setData({ sonTasks });
      }
    },

    // 撤销按钮事件
    onCancel() {
      this.triggerEvent('cancel');
    },

    // 应用按钮事件
    onConfirm() {
      wx.showModal({
        title: '确认创建',
        content: '确认创建计划,以及对应的子任务?',
        success: (res) => {
          if (res.confirm) {
            this.triggerEvent('confirm', this.properties.aiPlanText);
          }
        },
      });
    }
  }
});
