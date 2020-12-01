<template>
  <div>
  <el-row style="margin-bottom: 30px">
      <el-form :inline="true">
        <el-form-item v-for="opt in options" :key="opt.key" :label="opt.name">
          <el-select clearable v-model="opt.value">
          <el-option
            v-for="item in opt.option"
            :key="item.value"
            :label="item.value"
            :value="item.value">
            <div v-if="item.label[0]=='<'" v-html="item.label"/>
            <img v-else :src="item.label">
          </el-option>
        </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="updateOption">保存</el-button>
          <el-button type="primary" @click="download">下载</el-button>
        </el-form-item>
      </el-form>
  </el-row>
  <el-row style="margin-bottom: 50px">
    <el-upload
      class="upload-demo"
      drag multiple
      ref="upload"
      action="api/upload"
      accept="image/jpeg"
      list-type="picture-card"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-remove="handleRemove"
    >
      <i class="el-icon-upload">
        <div class="el-upload__text">
          将文件拖到此处，或
          <em>点击上传</em>
        </div>
        <div class="el-upload__tip" slot="tip">只能上传jpg文件</div>
      </i>
    </el-upload>
  </el-row>
  <el-divider></el-divider>
  <el-row v-for="(item, key, index) in svgs" :key="index">
    <el-form :inline="true">
      <el-form-item>
        <el-input style="" @change="changeFileName(key, item, ...arguments)" v-model="fileNames[index]"></el-input>
      </el-form-item>
      <el-form-item>
        <div v-html="item" style="display: inline-block;"></div>
      </el-form-item>
    </el-form>
  </el-row>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: "HelloWorld",
  created () {
    document.title = ''
  },
  data() {
    return {
      svgs: {
        // filename: svg file
      },
      jsons: { // 文件名保持一致
        // filename: origin json data
      },
      options: [
        {key: "manufacturer",   name: "厂商",   value: "", option: []},
        {key: "network",        name: "网口",   value: "", option: []},
        {key: "optical",        name: "光口",   value: "", option: []},
        {key: "indicatorlight", name: "指示灯", value: "", option: []},
        {key: "usb",            name: "USB",    value: "", option: []},
        {key: "serial",         name: "接口ID", value: "", option: []}
      ],
      fileNames: [] 
    }
  },
  mounted() {
    for(let index in this.options) {
      let opt = this.options[index]
      this.queryOption(opt.key).then( data => {
        opt.value = data.default
        opt.option = data.option
        this.$set(this.options, index, opt)
      })
    }
  },
  methods: {
    handleError(err) {
      console.error(err)
    },
    handleSuccess(response, file) {
      let fileName = file.name.replace(".jpg","")
      this.$set(this.jsons, fileName, response)
      let params = {}
      params[fileName] = response
      axios.post('/api/reloadsvg', params)
        .then(res => {
          this.$set(this.svgs, fileName, res.data[fileName] || `ERROR`)
          this.fileNames = Object.keys(this.svgs)
        }).catch(e => {
          console.error(e)
        })
    },
    handleRemove(file) {
      this.$delete(this.svgs, file.name.replace(".jpg",""))
      this.$delete(this.jsons, file.name.replace(".jpg",""))
    },
    updateOption() {
      let params = this.options.map(opt => {
        return {label:opt.key,value:opt.value}
      })
      axios.post('/api/option', params)
        .then(res => {
          if(res.data && res.data.status) {
            axios.post('/api/reloadsvg', this.jsons)
              .then(res => {
                for(let key in res.data) {
                  this.$set(this.svgs, key, res.data[key])
                }
              }).catch(e => {
                console.error(e)
        })
          }
        }).catch(e => {
          console.error(e)
        })
    },
    download() {
      this.postDownload('/api/download', this.jsons)
    },
    async queryOption(type) {
      let ret = []
      await axios.get('/api/option', {params:{key:type}})
        .then(res => {
          ret = res.data
        }).catch(e => {
          console.error(e)
        })
      return ret
    },
    changeFileName (key, item, val) {
      this.$delete(this.svgs, key)
      this.$set(this.jsons, val, this.jsons[key])
      this.$delete(this.jsons, key)
      this.$set(this.svgs, val, item)
    },
    async postDownload (url, params) {
      const request = {
        body: JSON.stringify(params),
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=UTF-8'
        }
      }
      const response = await fetch(url, request)
      const blob = await response.blob()
      const link = document.createElement('a')
      link.download = decodeURIComponent('test')
      link.style.display = 'none'
      link.href = URL.createObjectURL(blob)
      document.body.appendChild(link)
      link.click()
      URL.revokeObjectURL(link.href)
      document.body.removeChild(link)
      
    }
  }
};
</script>

<style scoped>
.el-upload__tip {
  margin-top: 0px;
}
.el-select-dropdown__item{
  background-color: silver
}
</style>>
<style>
  .el-upload--picture-card {
  width: 360px;
}
</style>