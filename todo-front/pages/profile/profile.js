import * as echarts from '../../components/ec-canvas/echarts';
let api = require('../../common/api');

let chart = null;

Page({
  data: {
    userInfo: {},
    taskCount: 0,
    completionRate: 0,
    levelTitle: "Lv.1初出茅庐",
    barChart: null,
    lineChart: null,
    // ecWordCloud: {
    //   lazyLoad: true
    // },
    // word_cloud: [] // 你的原始词云数据

    barXData: [],
    barYData: [],
    lineXData: [],
    lineYData: [],
  },

  async onLoad() {
    this.getUserInfo();
    this.initData();
    await this.initCharts();
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

  // 初始化数据
  initData() {
    wx.request({
      url: api.ApiHost + '/profile/data',
      method: 'get',
      data: {
        // TODO: 改为年获取
        "user_id": 1, 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          return 
        }

        this.setData({
          taskCount: res.data.total_task_len,
          completionRate: (res.data.task_finished_rate * 100).toFixed(1),
          levelTitle: this.initLevelTitle(res.data.total_task_len),
          barXData: res.data.bar_chart?.x_axis,
          barYData: res.data.bar_chart?.y_axis,
          lineXData: res.data.line_chart?.x_axis,
          lineYData: res.data.line_chart?.y_axis,
        })
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

  // 等级
  initLevelTitle(total_task_len) {
    if (total_task_len >= 1000) {
      return "Lv.4 炉火纯青"
    }
    if (total_task_len >= 100) {
      return "Lv.3 独当一面"
    }
    if (total_task_len >= 10) {
      return "Lv.2 崭露头角"
    }
    if (total_task_len >= 0) {
      return "Lv.1 初出茅庐"
    }
  },

  async initCharts() {
    // 初始化柱状图
    this.initBarChart();
    // 初始化折线图
    this.initLineChart();
    // TODO: 没搞出来
    // this.initWordCloud();
  },

  initBarChart() {
    this.selectComponent('#barChart').init((canvas, width, height, dpr) => {
      const chart = echarts.init(canvas, null, {
        width: width,
        height: height,
        devicePixelRatio: dpr
      });
      
      const option = {
        title: {
          text: '任务完成量',
          left: 'center',
          top: 10,
          textStyle: {
              fontSize: 12,
              color: '#666666',  // 更改为更浅的灰色
              fontWeight: 'rgba(102, 102, 102, 0.8)'  // 可以设置字体粗细
          }
        },
        xAxis: {
          type: 'category',
          data: this.data.barXData,
          axisLine: {
            lineStyle: {
              color: '#999'
            }
          }
        },
        yAxis: {
          type: 'value',
          splitLine: {
            lineStyle: {
              color: '#eee'
            }
          }
        },
        series: [{
          data: this.data.barYData,
          type: 'bar',
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0,
                color: '#FFB800' // 渐变开始颜色
              }, {
                offset: 1,
                color: '#FF9900' // 渐变结束颜色
              }]
            }
          },
          barWidth: '40%'
        }]
      };
      
      chart.setOption(option);
      return chart;
    });
  },

  // 折线图配置
  initLineChart() {
    this.selectComponent('#lineChart').init((canvas, width, height, dpr) => {
      const chart = echarts.init(canvas, null, {
        width: width,
        height: height,
        devicePixelRatio: dpr
      });
      
      const option = {
        title: {
          text: '任务完成率',
          left: 'center',
          top: 10,
          textStyle: {
              fontSize: 12,
              color: '#666666',  // 更改为更浅的灰色
              fontWeight: 'rgba(102, 102, 102, 0.8)'  // 可以设置字体粗细
          }
        },
        xAxis: {
          type: 'category',
          data: this.data.lineXData,
          axisLine: {
            lineStyle: {
              color: '#999'
            }
          }
        },
        yAxis: {
          type: 'value',
          splitLine: {
            lineStyle: {
              color: '#eee'
            }
          }
        },
        series: [{
          data: this.data.lineYData,
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 8,
          lineStyle: {
            color: '#FF9900',
            width: 3
          },
          itemStyle: {
            color: '#FF9900',
            borderWidth: 2,
            borderColor: '#fff'
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0,
                color: 'rgba(255, 153, 0, 0.3)' // 渐变开始颜色，较透明
              }, {
                offset: 1,
                color: 'rgba(255, 153, 0, 0.05)' // 渐变结束颜色，更透明
              }]
            }
          }
        }]
      };
      
      chart.setOption(option);
      return chart;
    });
  },

  // 初始化词云图
  initWordCloud() {
    wx.request({
      url: api.ApiHost + '/profile/data',
      method: 'get',
      data: {
        // TODO: 改为年获取
        "user_id": 1, 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'Bearer ' + api.AuthKey // 如果需要token
      },
      success: (res) => {
        if (!res || !res.data) {
          return 
        }

        this.setData({
          word_cloud: res.data.word_cloud
        })
        console.log(res.data.word_cloud)
        this.rendertWordCloud()
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

  // 渲染词云
  rendertWordCloud() {
    const component = this.selectComponent('#wordCloudChart');
    if (!component) {
      console.error('Cannot find component #wordCloudChart');
      return;
    }

    component.init((canvas, width, height, dpr) => {
      chart = echarts.init(canvas, null, {
        width: width,
        height: height,
        devicePixelRatio: dpr
      });
      canvas.setChart(chart);

      // 转换数据
      const transformedData = this.transformWordCloudData(this.data.word_cloud, width, height);

      const option = {
        grid: {
          top: 10,
          bottom: 10,
          left: 10,
          right: 10,
          containLabel: true
        },
        xAxis: {
          type: 'value',
          show: false,
          min: 0,
          max: width
        },
        yAxis: {
          type: 'value',
          show: false,
          min: 0,
          max: height
        },
        series: [{
          type: 'scatter',
          data: transformedData,
          symbolSize: 1,
          label: {
            show: true,
            formatter: function(param) {
              return param.data[2];
            },
            position: 'inside',
            fontSize: function(param) {
              // 根据value值计算字体大小，范围在12-30之间
              const value = param.data[3];
              return Math.max(12, Math.min(30, value * 8 + 12));
            },
            color: function() {
              // 随机生成暖色调
              return 'rgb(' + [
                Math.round(Math.random() * 50) + 200,
                Math.round(Math.random() * 100) + 100,
                Math.round(Math.random() * 50)
              ].join(',') + ')';
            },
            fontWeight: 'bold',
            textShadowBlur: 5,
            textShadowColor: 'rgba(0, 0, 0, 0.3)'
          },
          itemStyle: {
            color: 'transparent'
          },
          animation: true,
          animationDuration: 1000,
          animationDelay: function(idx) {
            return idx * 100;
          }
        }]
      };

      chart.setOption(option);
      return chart;
    });
  },

  // 转换数据格式的函数
  transformWordCloudData(wordCloud, width, height) {
    // 找出最大和最小值，用于计算字体大小的比例
    const values = wordCloud.map(item => item.value);
    const maxValue = Math.max(...values);
    const minValue = Math.min(...values);

    // 计算位置的辅助变量
    const centerX = width / 2;
    const centerY = height / 2;
    const radius = Math.min(width, height) / 3;

    return wordCloud.map((item, index) => {
      // 使用螺旋形布局计算位置
      const angle = (index / wordCloud.length) * Math.PI * 8; // 8表示螺旋圈数
      const spiralRadius = (radius * (index + 1)) / wordCloud.length;
      const x = centerX + Math.cos(angle) * spiralRadius;
      const y = centerY + Math.sin(angle) * spiralRadius;

      return [x, y, item.name, item.value];
    });
  },

  // 页面相关事件处理函数
  onPullDownRefresh() {
    // 下拉刷新逻辑
    this.getUserInfo();
    this.initCharts();
    wx.stopPullDownRefresh();
  },

  onShareAppMessage() {
    return {
      title: '我的个人数据报告',
      path: '/pages/profile/profile'
    };
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