<template>
  <div>
  <el-row>
    <el-upload
      class="upload-demo"
      drag multiple
      ref="upload"
      action="/upload"
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
  <el-row style="margin-top: 30px">
    <div v-for="item, key, index in svgs">
      <div v-html="item"></div>
    </div>
  </el-row>
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  data() {
    return {
      svgs: {}
    }
  },
  methods: {
    handleError(err) {
      console.log(err)
    },
    handleSuccess(response, file) {
      this.$set(this.svgs, file.name, response)
    },
    handleRemove(file) {
      this.$delete(this.svgs, file.name)
    }
  }
};
</script>

<style scoped>
.el-upload__tip {
  margin-top: 0px;
}
</style>>