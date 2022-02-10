import request from '@/utils/request'

export function serviceList(query) {
  return request({
    url: '/service/serviceList',
    method: 'get',
    params: query
  })
}

export function serviceDelete(deleteQuery) {
  return request({
    url: '/service/serviceDelete',
    method: 'get',
    params: deleteQuery
  })
}

export function serviceDetail(query) {
  return request({
    url: '/service/serviceDetail',
    method: 'get',
    params: query
  })
}

export function serviceAddHttp(query) {
  return request({
    url: '/service/serviceAddHttp',
    method: 'post',
    params: query
  })
}

export function serviceUpdateHttp(query) {
  return request({
    url: '/service/serviceUpdateHttp',
    method: 'post',
    params: query
  })
}

export function serviceAddTcp(query) {
  return request({
    url: '/service/serviceAddTcp',
    method: 'post',
    params: query
  })
}

export function serviceUpdateTcp(query) {
  return request({
    url: '/service/serviceUpdateTcp',
    method: 'post',
    params: query
  })
}

export function serviceAddGrpc(query) {
  return request({
    url: '/service/serviceAddGrpc',
    method: 'post',
    params: query
  })
}

export function serviceUpdateGrpc(query) {
  return request({
    url: '/service/serviceUpdateGrpc',
    method: 'post',
    params: query
  })
}

export function serviceStat(query) {
  return request({
    url: '/service/serviceStat',
    method: 'get',
    params: query
  })
}
