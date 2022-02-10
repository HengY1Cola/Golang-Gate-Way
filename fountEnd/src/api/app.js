import request from '@/utils/request'

export function appList(query) {
  return request({
    url: '/app/appList',
    method: 'get',
    params: query
  })
}

export function appDelete(query) {
  return request({
    url: '/app/appDelete',
    method: 'get',
    params: query
  })
}

export function appAdd(query) {
  return request({
    url: '/app/appAdd',
    method: 'post',
    params: query
  })
}

export function appUpdate(query) {
  return request({
    url: '/app/appUpdate',
    method: 'post',
    params: query
  })
}

export function appDetail(query) {
  return request({
    url: '/app/appDetail',
    method: 'get',
    params: query
  })
}

export function appStat(query) {
  return request({
    url: '/app/appStat',
    method: 'get',
    params: query
  })
}
