<template>
  <div class="mixin-components-container">
    <el-row>
      <el-form ref="form" :model="form" label-width="130px">
        <el-card class="box-card" style="margin-bottom:20px;">
          <div slot="header" class="clearfix">
            <span v-if="is_edit === false">创建GRPC服务</span>
            <span v-if="is_edit === true">修改GRPC服务</span>
          </div>
          <div>
            <el-form-item label="服务名称" class="is-required">
              <el-input v-model="form.serviceName" :disabled="is_edit===true" placeholder="6-128位字母数字下划线" />
            </el-form-item>
            <el-form-item label="服务描述" class="is-required">
              <el-input v-model="form.serviceDesc" placeholder="最多255字符，必填" />
            </el-form-item>
            <el-form-item label="端口" class="is-required">
              <el-input v-model="form.port" :disabled="is_edit===true" placeholder="需要设置8001-8999范围内数字，必填" />
            </el-form-item>
            <el-form-item label="metadata转换">
              <el-input v-model="form.headerTransfor" placeholder="metadata转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多条换行" type="textarea" :autosize="{ minRows: 2, maxRows: 20}" />
            </el-form-item>
            <el-form-item label="开启验证">
              <el-switch
                v-model="form.openAuth"
                active-value="1"
                inactive-value="0"
              />
            </el-form-item>
            <el-form-item label="IP白名单">
              <el-input v-model="form.whiteList" placeholder="格式: 127.0.0.1 多条换行，白名单优先级高于黑名单" type="textarea" :autosize="{ minRows: 2, maxRows: 20}" />
            </el-form-item>

            <el-form-item label="IP黑名单">
              <el-input v-model="form.blackList" placeholder="格式: 127.0.0.1 多条换行" type="textarea" :autosize="{ minRows: 2, maxRows: 20}" />
            </el-form-item>

            <el-form-item label="客户端限流">
              <el-input v-model="form.clientIpFlowLimit" placeholder="0表示无限制" />
            </el-form-item>

            <el-form-item label="服务端限流">
              <el-input v-model="form.serviceIpFlowLimit" placeholder="0表示无限制" />
            </el-form-item>

            <el-form-item label="轮询方式">
              <el-radio v-model="form.roundType" label="0">random</el-radio>
              <el-radio v-model="form.roundType" label="1">round-robin</el-radio>
              <el-radio v-model="form.roundType" label="2">weight_round-robin</el-radio>
              <el-radio v-model="form.roundType" label="3">ip_hash</el-radio>
            </el-form-item>

            <el-form-item label="IP列表" class="is-required">
              <el-input v-model="form.ipList" placeholder="格式: 127.0.0.1:80 多条换行" type="textarea" :autosize="{ minRows: 2, maxRows: 20}" />
            </el-form-item>

            <el-form-item label="权重列表" class="is-required">
              <el-input v-model="form.weightList" placeholder="格式: 50 多条换行" type="textarea" :autosize="{ minRows: 2, maxRows: 20}" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="disableButton" @click="onSubmit">立即提交</el-button>
            </el-form-item>
          </div>
        </el-card>
      </el-form>
    </el-row>
  </div>
</template>

<script>
import { serviceDetail, serviceAddGrpc, serviceUpdateGrpc } from '@/api/service'
export default {
  components: { },
  data() {
    return {
      is_edit: false,
      disableButton: false,
      form: {
        id: '0',
        serviceName: '',
        serviceDesc: '',
        port: '',
        ruleType: '0',
        rule: '',
        needHttps: '0',
        needStripUri: '1',
        needWebsocket: '1',
        urlRewrite: '',
        headerTransfor: '',
        openAuth: '0',
        blackList: '',
        whiteList: '',
        whiteHostName: '',
        clientIpFlowLimit: '',
        serviceIpFlowLimit: '',
        roundType: '2',
        ipList: '',
        weightList: '',
        forbidList: '',
        upstreamConnectTimeout: '',
        upstreamHeaderTimeout: '',
        upstreamIdleTimeout: '',
        upstreamMaxIdle: ''
      }
    }
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    if (id > 0) {
      this.is_edit = true
      this.fetchData(id)
    }
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    fetchData(id) {
      this.listLoading = true
      const query = {
        id: id
      }
      serviceDetail(query).then(response => {
        const formData = response.data
        this.form.id = formData.info.id
        this.form.serviceName = formData.info.service_name
        this.form.serviceDesc = formData.info.service_desc
        this.form.port = formData.grpc.port.toString()
        this.form.ruleType = formData.http.rule_type.toString()
        this.form.rule = formData.http.rule
        this.form.needHttps = formData.http.need_https.toString()
        this.form.needStripUri = formData.http.need_strip_uri.toString()
        this.form.needWebsocket = formData.http.need_websocket.toString()
        this.form.urlRewrite = formData.http.url_rewrite.replace(/,/g, '\n')
        this.form.headerTransfor = formData.grpc.header_transfor.replace(/,/g, '\n')
        this.form.openAuth = formData.accessControl.open_auth.toString()
        this.form.blackList = formData.accessControl.black_list.replace(/,/g, '\n')
        this.form.whiteList = formData.accessControl.white_list.replace(/,/g, '\n')
        this.form.whiteHostName = formData.accessControl.white_host_name.replace(/,/g, '\n')
        this.form.clientIpFlowLimit = formData.accessControl.clientip_flow_limit
        this.form.serviceIpFlowLimit = formData.accessControl.service_flow_limit
        this.form.roundType = formData.loadBalance.round_type.toString()
        this.form.ipList = formData.loadBalance.ip_list.replace(/,/g, '\n')
        this.form.weightList = formData.loadBalance.weight_list.replace(/,/g, '\n')
        this.form.forbidList = formData.loadBalance.forbid_list.replace(/,/g, '\n')
        this.form.upstreamConnectTimeout = formData.loadBalance.upstream_connect_timeout
        this.form.upstreamHeaderTimeout = formData.loadBalance.upstream_header_timeout
        this.form.upstreamIdleTimeout = formData.loadBalance.upstream_idle_timeout
        this.form.upstreamMaxIdle = formData.loadBalance.upstream_max_idle
        this.listLoading = false
      })
    },
    onSubmit() {
      this.disableButton = true
      const formData = Object.assign({}, this.form)
      formData.port = Number(formData.port)
      formData.ruleType = Number(formData.ruleType)
      formData.needHttps = Number(formData.needHttps)
      formData.needStripUri = Number(formData.needStripUri)
      formData.needWebsocket = Number(formData.needWebsocket)
      formData.openAuth = Number(formData.openAuth)
      formData.roundType = Number(formData.roundType)
      formData.clientIpFlowLimit = Number(formData.clientIpFlowLimit)
      formData.serviceIpFlowLimit = Number(formData.serviceIpFlowLimit)
      formData.upstreamConnectTimeout = Number(formData.upstreamConnectTimeout)
      formData.upstreamHeaderTimeout = Number(formData.upstreamHeaderTimeout)
      formData.upstreamIdleTimeout = Number(formData.upstreamIdleTimeout)
      formData.upstreamMaxIdle = Number(formData.upstreamMaxIdle)
      formData.ipList = formData.ipList.replace(/\n/g, ',')
      formData.weightList = formData.weightList.replace(/\n/g, ',')
      formData.forbidList = formData.forbidList.replace(/\n/g, ',')
      formData.whiteHostName = formData.whiteHostName.replace(/\n/g, ',')
      formData.whiteList = formData.whiteList.replace(/\n/g, ',')
      formData.blackList = formData.blackList.replace(/\n/g, ',')
      formData.headerTransfor = formData.headerTransfor.replace(/\n/g, ',')
      formData.urlRewrite = formData.urlRewrite.replace(/\n/g, ',')
      if (this.is_edit) {
        serviceUpdateGrpc(formData).then(() => {
          this.disableButton = false
          this.$notify({
            title: 'Success',
            message: '更改成功',
            type: 'success',
            duration: 1500
          })
        }, () => {
          this.disableButton = false
        })
      } else {
        serviceAddGrpc(formData).then(() => {
          this.disableButton = false
          this.$notify({
            title: 'Success',
            message: '添加成功',
            type: 'success',
            duration: 1500
          })
        }, () => {
          this.disableButton = false
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
.component-item{
  min-height: 100px;
}
.el-select .el-input {
  width: 130px;
}
.input-with-select .el-input-group__prepend {
  background-color: #fff;
}
</style>
