Component({
  properties: {
    aiFormatText: {
      type: String,
      observer: function(newVal) {
        if (newVal) {
          this.excuteParse()
        }
      }
    }, 
    // xx日期下的所有未完成待办列表
    dateTodoList: {
      type: Array,
      value: null,
      observer: function(newVal) {
        if (newVal) {
          this.excuteParse()
        }
      }
    }, 
  },

  data: {
    tasks: [] // 存储所有解析后的任务
  },

  methods: {
    excuteParse() {
      let aiFormatText = this.data.aiFormatText
      let dateTodoList = this.data.dateTodoList
      if (aiFormatText != "" && dateTodoList && dateTodoList.length > 0) {
        this.parseAiContent()
      }
    },

    parseAiContent() {
      // 按换行符分割不同的操作
      const operations = this.data.aiFormatText.split('\n').filter(op => op.trim());
      const parsedTasks = [];

      operations.forEach(operation => {
        if (operation.includes('[[SplitTask]]')) {
          parsedTasks.push(this.parseSplitTask(operation));
        } else if (operation.includes('[[UpdateNameTask]]')) {
          parsedTasks.push(this.parseUpdateNameTask(operation));
        } else if (operation.includes('[[UpdateDateTask]]')) {
          parsedTasks.push(this.parseUpdateDateTask(operation));
        }
      });

      this.setData({ tasks: parsedTasks });
    },

    parseSplitTask(content) {
      const parts = content.split('|||');
      const deleteTask = parts[0].match(/\[delete\](\d+)/)[1];
      const addTasks = parts.slice(1).map(task => {
        const match = task.match(/\[add\](.*?)@(.*)/);
        return {
          name: match[1],
          date: match[2]
        };
      });

      let taskDetail = this.getTaskDetail(deleteTask);
      return {
        type: 'split',
        originalTask: {
          id: taskDetail?.id,
          name: taskDetail?.content, 
          date: taskDetail?.date,
        },
        newTasks: addTasks
      };
    },

    parseUpdateNameTask(content) {
      const match = content.match(/\[update_name\](\d+)@(.*)/);
      let taskDetail = this.getTaskDetail(match[1]);
      return {
        type: 'updateName',
        originalTask: {
          id: taskDetail?.id,
          name: taskDetail?.content, 
          date: taskDetail?.date,
        },
        newName: match[2]
      };
    },

    parseUpdateDateTask(content) {
      const match = content.match(/\[update_date\](\d+)@(.*)/);
      let taskDetail = this.getTaskDetail(match[1]);
      return {
        type: 'updateDate',
        originalTask: {
          id: taskDetail?.id,
          name: taskDetail?.content, 
          date: taskDetail?.date,
        },
        newDate: match[2]
      };
    }, 

    // 通过任务id获取任务详情
    getTaskDetail(taskID) {
      for(let item of this.data.dateTodoList) {
        if (item.id == taskID) {
          return item
        } 
      }
      return null
    }
  }
});
