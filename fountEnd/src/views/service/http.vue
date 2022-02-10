<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span v-if="isEdit===false">创建HTTP服务</span>
          <span v-if="isEdit===true">修改HTTP服务</span>
        </div>
        <div style="margin-bottom:50px;">
          <el-form ref="form" :model="form" label-width="140px">
            <el-form-item label="服务名称" class="is-required">
              <el-input v-model="form.serviceName" placeholder="6-128位字母数字下划线" :disabled="isEdit===true" />
            </el-form-item>
            <el-form-item label="服务描述" class="is-required">
              <el-input v-model="form.serviceDesc" placeholder="最多255字符，必填" />
            </el-form-item>
            <el-form-item label="接入类型" class="is-required">
              <el-input v-model="form.rule" :disabled="isEdit===true" placeholder="路径格式：/user/,域名格式：www.test.com" class="input-with-select">
                <el-select slot="prepend" v-model="form.ruleType" placeholder="请选择" style="width:80px" :disabled="isEdit===true">
                  <el-option label="路径" :value="0" />
                  <el-option label="域名" :value="1" />
                </el-select>
              </el-input>
            </el-form-item>
            <el-form-item label="支持https">
              <el-switch
                v-model="form.needHttps"
                :active-value="1"
                :inactive-value="0"
              />
              <span style="color:#606266;font-weight: 700; margin-left: 20px">支持strip_uri</span>
              <el-switch
                v-model="form.needStripUrl"
                :active-value="1"
                :inactive-value="0"
                style="margin-left: 10px"
              />
              <span style="color:#606266;font-weight: 700;margin-left: 20px">支持websocket</span>
              <el-switch
                v-model="form.needWebsocket"
                :active-value="1"
                :inactive-value="0"
                style="margin-left: 10px"
              />
            </el-form-item>
            <el-form-item label="URL重写">
              <el-input v-model="form.urlRewrite" type="textarea" autosize placeholder="格式：^/gateway/test_service(.*) $1 多条换行" />
            </el-form-item>
            <el-form-item label="Header转换">
              <el-input v-model="form.headerTransfor" type="textarea" autosize placeholder="header转换支持 add(增加)/del(删除)/edit(修改) 格式：add headerName headValue" />
            </el-form-item>
            <el-form-item label="开启验证">
              <el-switch
                v-model="form.openAuth"
                :active-value="1"
                :inactive-value="0"
              />
            </el-form-item>
            <el-form-item label="IP白名单">
              <el-input v-model="form.whiteList" type="textarea" autosize placeholder="格式：127.0.0.1 多条换行，白名单优先于黑名单" />
            </el-form-item>
            <el-form-item label="IP黑名单">
              <el-input v-model="form.blackList" type="textarea" autosize placeholder="格式：127.0.0.1 多条换行，白名单优先于黑名单" />
            </el-form-item>
            <el-form-item label="客户端限流">
              <el-input v-model="form.clientIpFlowLimit" placeholder="0表示无限制" />
            </el-form-item>
            <el-form-item label="服务端限流">
              <el-input v-model="form.serviceIpFlowLimit" placeholder="0表示无限制" />
            </el-form-item>
            <el-form-item label="轮询方式">
              <el-radio-group v-model="form.roundType">
                <el-radio :label="0">random</el-radio>
                <el-radio :label="1">round-robin</el-radio>
                <el-radio :label="2">weight_round_robin</el-radio>
                <el-radio :label="3">ip_hash</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="IP列表" class="is-required">
              <el-input v-model="form.ipList" type="textarea" autosize placeholder="格式：127.0.0.1:80 多条换行" />
            </el-form-item>
            <el-form-item label="权重列表" class="is-required">
              <el-input v-model="form.weightList" type="textarea" autosize placeholder="格式：50 多条换行" />
            </el-form-item>
            <el-form-item label="建立连接超时">
              <el-input v-model="form.upstreamConnectTimeout" placeholder="单位s，0表示无限制" />
            </el-form-item>
            <el-form-item label="获取header超时">
              <el-input v-model="form.upstreamHeaderTimeout" placeholder="单位s，0表示无限制" />
            </el-form-item>
            <el-form-item label="链接最大空闲时间">
              <el-input v-model="form.upstreamIdleTimeout" placeholder="单位s，0表示无限制" />
            </el-form-item>
            <el-form-item label="最大空闲链接数">
              <el-input v-model="form.upstreamMaxIdle" placeholder="0表示无限制" />
            </el-form-item>
            <!--            提交按钮-->
            <el-form-item>
              <el-button type="primary" :disabled="submitButDisabled" @click="handleSubmit">立即提交</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import { serviceAddHttp, serviceDetail, serviceUpdateHttp } from '@/api/service'
export default {
  name: 'ServiceCreateHttp',
  data() {
    return {
      isEdit: false,
      submitButDisabled: false,
      form: {
        serviceName: '',
        serviceDesc: '',
        ruleType: 0,
        rule: '',
        needHttps: 0,
        needWebsocket: 0,
        needStripUrl: 0,
        urlRewrite: '',
        headerTransfor: '',
        roundType: 2,
        ipList: '',
        weightList: '',
        upstreamConnectTimeout: '',
        upstreamHeaderTimeout: '',
        upstreamIdleTimeout: '',
        upstreamMaxIdle: '',
        openAuth: 0,
        blackList: '',
        whiteList: '',
        clientIpFlowLimit: '',
        serviceIpFlowLimit: ''
      }
    }
  },
  created() { // 判断是创建还是修改
    const id = this.$route.params && this.$route.params.id
    if (id > 0) {
      this.isEdit = true
      this.fetchData(id)
    }
  },
  methods: {
    fetchData(id) { // 拿到细节内容
      const query = { 'id': id }
      serviceDetail(query).then(response => {
        this.form.id = response.data.info.id
        this.form.load_type = response.data.info.load_type
        this.form.serviceName = response.data.info.service_name
        this.form.serviceDesc = response.data.info.service_desc
        this.form.ruleType = response.data.http.rule_type
        this.form.rule = response.data.http.rule
        this.form.needHttps = response.data.http.need_https
        this.form.needWebsocket = response.data.http.need_websocket
        this.form.needStripUrl = response.data.http.need_strip_uri
        this.form.urlRewrite = response.data.http.url_rewrite.replace(/,/g, '\n')
        this.form.headerTransfor = response.data.http.header_transfor.replace(/,/g, '\n')
        this.form.roundType = response.data.loadBalance.round_type
        this.form.ipList = response.data.loadBalance.ip_list.replace(/,/g, '\n')
        this.form.weightList = response.data.loadBalance.weight_list.replace(/,/g, '\n')
        this.form.upstreamConnectTimeout = response.data.loadBalance.upstream_connect_timeout
        this.form.upstreamHeaderTimeout = response.data.loadBalance.upstream_header_timeout
        this.form.upstreamIdleTimeout = response.data.loadBalance.upstream_idle_timeout
        this.form.upstreamMaxIdle = response.data.loadBalance.upstream_max_idle
        this.form.openAuth = response.data.accessControl.open_auth
        this.form.blackList = response.data.accessControl.black_list.replace(/,/g, '\n')
        this.form.whiteList = response.data.accessControl.white_list.replace(/,/g, '\n')
        this.form.clientIpFlowLimit = response.data.accessControl.clientip_flow_limit
        this.form.serviceIpFlowLimit = response.data.accessControl.service_flow_limit
      }).catch(() => {
      })
    },
    handleSubmit() { // 构建提交表单
      this.submitButDisabled = true
      const query = Object.assign({}, this.form)
      query.urlRewrite = query.urlRewrite.replace(/\n/g, ',')
      query.headerTransfor = query.headerTransfor.replace(/\n/g, ',')
      query.ipList = query.ipList.replace(/\n/g, ',')
      query.weightList = query.weightList.replace(/\n/g, ',')
      query.blackList = query.blackList.replace(/\n/g, ',')
      query.whiteList = query.whiteList.replace(/\n/g, ',')
      query.serviceIpFlowLimit = Number(query.serviceIpFlowLimit)
      query.needStripUrl = Number(query.needStripUrl)
      query.needHttps = Number(query.needHttps)
      query.needWebsocket = Number(query.needWebsocket)
      query.clientIpFlowLimit = Number(query.clientIpFlowLimit)
      query.upstreamConnectTimeout = Number(query.upstreamConnectTimeout)
      query.upstreamHeaderTimeout = Number(query.upstreamHeaderTimeout)
      query.upstreamIdleTimeout = Number(query.upstreamIdleTimeout)
      query.upstreamMaxIdle = Number(query.upstreamMaxIdle)
      if (this.isEdit) {
        serviceUpdateHttp(query).then(response => {
          this.submitButDisabled = false
          this.$notify({
            title: 'Success',
            message: '修改成功',
            type: 'success',
            duration: 2000
          })
          // this.fetchData(this.$route.params && this.$route.params.id)
        }).catch(() => {
          this.submitButDisabled = false
        })
      } else {
        serviceAddHttp(query).then(response => {
          this.submitButDisabled = false
          this.$notify({
            title: 'Success',
            message: '添加成功',
            type: 'success',
            duration: 2000
          })
        }).catch(() => {
          this.submitButDisabled = false
        })
      }
    }
  }
}
</script>

<style scoped>
.mixin-components-container {
  background-color: #f0f2f5;
  padding: 30px;
  min-height: calc(100vh - 84px);
}
</style>
