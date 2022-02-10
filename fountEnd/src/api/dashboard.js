import request from '@/utils/request'

export function panelGroupData() {
  return request({
    url: '/main/panelGroupData',
    method: 'get'
  })
}

export function flowStat() {
  return request({
    url: '/main/flowStat  ',
    method: 'get'
  })
}

export function serviceStat() {
  return request({
    url: '/main/serviceStat',
    method: 'get'
  })
}
