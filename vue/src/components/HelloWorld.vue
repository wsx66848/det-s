<template>
  <div>
  <el-row style="margin-bottom: 30px">
    <el-col :span="4" v-for="opt in options" :key="opt.key">
      <el-select v-model="opt.value" :placeholder="opt.name">
        <el-option
          v-for="item in opt.option"
          :key="item.value"
          :label="item.value"
          :value="item.value">
          <div v-if="item.label[0]=='<'" v-html="item.label"/>
          <img v-else :src="item.label">
        </el-option>
      </el-select>
    </el-col>
    <el-col :span="4" style="white-space:nowrap;">
      <el-button type="primary" icon="el-icon-setting" @click="updateOption"/>
      <el-button type="primary" icon="el-icon-download" @click="download"/>
    </el-col>
  </el-row>
  <el-row style="margin-bottom: 30px">
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
  <el-row>
    <div v-for="item, key, index in svgs">
      <div v-html="item"></div>
    </div>
  </el-row>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: "HelloWorld",
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
        {key: "usb",            name: "USB",    value: "", option: []}
      ]
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
      this.$set(this.jsons, file.name, response)
      let params = {}
      params[file.name] = response
      axios.post('/api/reloadsvg', params)
        .then(res => {
          this.$set(this.svgs, file.name, res.data[file.name] || `ERROR`)
        }).catch(e => {
          console.error(e)
        })
    },
    handleRemove(file) {
      this.$delete(this.svgs, file.name)
      this.$delete(this.jsons, file.name)
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
    download() {},
    async queryOption(type) {
      let ret = []
      await axios.get('/api/option', {params:{key:type}})
        .then(res => {
          ret = res.data
        }).catch(e => {
          console.error(e)
        })
      return ret
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