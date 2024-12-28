package api

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
	"todo-ai/core/logger"
	"todo-ai/model"

	"github.com/gin-gonic/gin"
	"github.com/yanyiwu/gojieba"
)

type ProfileDataRequest struct {
	UserID int64 `json:"user_id" form:"user_id"` // 用户ID
}

type ProfileDataResp struct {
	TotalTaskLen     int64          `json:"total_task_len"`     // 总任务数量
	TaskFinishedRate float64        `json:"task_finished_rate"` // 任务完成率
	WordCloud        []*Word        `json:"word_cloud"`         // 词云
	BarChart         *BarChartData  `json:"bar_chart"`          // bar图表
	LineChart        *LineChartData `json:"line_chart"`         // line图表
}

type Word struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type BarChartData struct {
	XAxis []string `json:"x_axis"`
	YAxis []int    `json:"y_axis"`
}

type LineChartData struct {
	XAxis []string `json:"x_axis"`
	YAxis []int    `json:"y_axis"`
}

// ProfileData 个人信息数据
func ProfileData(ctx *gin.Context) (interface{}, error) {
	var request ProfileDataRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ProfileData BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	// 获取用户总任务数量
	totalTaskLen, err := model.GetTotalTaskLenByUserID(request.UserID)
	if err != nil {
		logger.Errorf("user_id[%d], model.GetTotalTaskLenByUserID error:%s", request.UserID, err)
		return nil, err
	}

	// 获取用户任务完成率
	taskFinishedRate, err := model.GetFinishedTaskLenByUserID(request.UserID)
	if err != nil {
		logger.Errorf("user_id[%d], model.GetFinishedTaskLenByUserID error:%s", request.UserID, err)
		return nil, err
	}

	// 获取词云数据
	wordCloud, err := ProfileWordCloud(request.UserID)
	if err != nil {
		logger.Errorf("user_id[%d], ProfileWordCloud error:%s", request.UserID, err)
		return nil, err
	}

	// 获取bar图表数据
	barChartData, lineChartData, err := GetChartData(request.UserID)
	if err != nil {
		logger.Errorf("user_id[%d], GetChartData error:%s", request.UserID, err)
		return nil, err
	}

	resp := &ProfileDataResp{
		TotalTaskLen:     totalTaskLen,
		TaskFinishedRate: float64(taskFinishedRate) / float64(totalTaskLen),
		WordCloud:        wordCloud,
		BarChart:         barChartData,
		LineChart:        lineChartData,
	}

	return resp, nil
}

// 获取词云数据
func ProfileWordCloud(userID int64) ([]*Word, error) {
	// 获取用户近30天的任务数据
	date := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	recentlyTasks, err := model.GetTaskByUserIDAndRecentlyDate(userID, date)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndRecentlyDate error:%s", err)
		return nil, err
	}

	recentTasksName := make([]string, 0)
	for _, task := range recentlyTasks {
		recentTasksName = append(recentTasksName, task.Name)
	}
	// 获取词云数据
	keywords := GetAllKeywords(recentTasksName)

	// 格式转化
	words := make([]*Word, 0)
	for k, v := range keywords {
		words = append(words, &Word{
			Name:  k,
			Value: v,
		})
	}

	return words, nil
}

// 获取图表数据
func GetChartData(userID int64) (*BarChartData, *LineChartData, error) {
	// 获取从今天开始往前7天的日期
	today := time.Now()
	sevenDaysAgo := today.AddDate(0, 0, -6) // -6是因为包含今天在内的7天
	startDate := sevenDaysAgo.Format("2006-01-02")

	// 获取用户近7天的任务数据
	recentlyTasks, err := model.GetTaskByUserIDAndRecentlyDate(userID, startDate)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndRecentlyDate error:%s", err)
		return nil, nil, err
	}

	// 准备日期列表（从7天前到今天）
	dates := make([]string, 7)
	for i := 0; i < 7; i++ {
		dates[i] = sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
	}

	barChartMap := make(map[string]int)  // 日期 -> 任务数量
	completedMap := make(map[string]int) // 日期 -> 已完成任务数量
	lineChartMap := make(map[string]int) // 日期 -> 任务完成率

	// 初始化所有日期的默认值
	for _, date := range dates {
		barChartMap[date] = 0
		completedMap[date] = 0
		lineChartMap[date] = 0
	}

	// 统计每天的任务总数和完成数
	for _, task := range recentlyTasks {
		if _, exists := barChartMap[task.Date]; exists { // 只统计我们关心的7天内的数据
			barChartMap[task.Date]++
			if task.Progress == 1 {
				completedMap[task.Date]++
			}
		}
	}

	// 计算每天的完成率
	for date, totalTasks := range barChartMap {
		if totalTasks > 0 {
			completionRate := int(math.Round(float64(completedMap[date]) * 100 / float64(totalTasks)))
			lineChartMap[date] = completionRate
		}
	}

	return ConvertMapToBarChart(barChartMap, dates), ConvertMapToLineChart(lineChartMap, dates), nil
}

// FormatDateToMonthDay 将 "2024-01-02" 格式转换为 "1.2" 格式
func FormatDateToMonthDay(dateStr string) string {
	// 解析原始日期
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr // 如果解析失败，返回原始字符串
	}
	// 格式化为 "月.日" 格式
	return fmt.Sprintf("%d.%d", t.Month(), t.Day())
}

// ConvertMapToLineChart 转换为折线图数据
func ConvertMapToLineChart(lineChartMap map[string]int, dates []string) *LineChartData {
	result := &LineChartData{
		XAxis: make([]string, len(dates)),
		YAxis: make([]int, len(dates)),
	}

	// 转换日期格式并填充数据
	for i, date := range dates {
		result.XAxis[i] = FormatDateToMonthDay(date)
		result.YAxis[i] = lineChartMap[date]
	}

	return result
}

// ConvertMapToBarChart 转换为柱状图数据
func ConvertMapToBarChart(barChartMap map[string]int, dates []string) *BarChartData {
	result := &BarChartData{
		XAxis: make([]string, len(dates)),
		YAxis: make([]int, len(dates)),
	}

	// 转换日期格式并填充数据
	for i, date := range dates {
		result.XAxis[i] = FormatDateToMonthDay(date)
		result.YAxis[i] = barChartMap[date]
	}

	return result
}

var (
	jieba *gojieba.Jieba
	once  sync.Once
)

func getJiebaInstance() *gojieba.Jieba {
	once.Do(func() {
		jieba = gojieba.NewJieba(
			"./dict/jieba.dict.utf8",
			"./dict/hmm_model.utf8",
			"./dict/user.dict.utf8",
			"./dict/idf.utf8",
			"./dict/stop_words.utf8",
		)
	})
	return jieba
}

func GetAllKeywords(dataList []string) map[string]int {
	keywords := make(map[string]int)
	jiebaInstance := getJiebaInstance()

	for _, danmu := range dataList {
		kds := ExtractKeywords(jiebaInstance, danmu)
		for _, kd := range kds {
			keywords[kd]++
		}
	}
	return keywords
}

func ExtractKeywords(jiebaInstance *gojieba.Jieba, text string) []string {
	keywords := make([]string, 0)

	words := jiebaInstance.ExtractWithWeight(text, 3) // 提取前3个关键词及其权重
	for _, word := range words {
		keywords = append(keywords, word.Word)
	}

	return keywords
}
