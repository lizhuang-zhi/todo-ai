// profile.js
import * as echarts from '../../components/ec-canvas/echarts';

Page({
  data: {
    userInfo: {},
    taskCount: 263,
    completionRate: 75,
    barChart: null,
    lineChart: null
  },

  onLoad() {
    this.getUserInfo();
    this.initCharts();
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

  initCharts() {
    // 初始化柱状图
    this.initBarChart();
    // 初始化折线图
    this.initLineChart();
  },

  initBarChart() {
    this.selectComponent('#barChart').init((canvas, width, height, dpr) => {
      const chart = echarts.init(canvas, null, {
        width: width,
        height: height,
        devicePixelRatio: dpr
      });
      
      const option = {
        xAxis: {
          type: 'category',
          data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
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
          data: [120, 200, 150, 80, 70],
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
        xAxis: {
          type: 'category',
          data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
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
          data: [820, 932, 901, 934, 1290, 1330, 1320],
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

  // 生成词云数据
  getWordCloudData() {
    return [
      { name: '任务A', value: 100 },
      { name: '项目B', value: 80 },
      { name: '工作C', value: 60 },
      { name: '计划D', value: 40 },
      { name: '目标E', value: 30 }
    ];
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
  }
});