const formatTime = date => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return `${[year, month, day].map(formatNumber).join('/')} ${[hour, minute, second].map(formatNumber).join(':')}`
}

const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : `0${n}`
}

// 将“Tue Dec 24 2024 22:15:08 GMT+0800 (中国标准时间)”转为“yyyy-MM-DD”
const formatDate = (date) => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0'); // 月份从0开始，需要+1
  const day = String(d.getDate()).padStart(2, '0');
  
  return `${year}-${month}-${day}`;
}

// 获取字符串的前n个字符，超过部分用...代替
const getShortStr = (str, n) => {
  return str.length > n ? str.slice(0, n) + '...' : str;
}

module.exports = {
  formatTime,
  formatDate,
  getShortStr,
}
