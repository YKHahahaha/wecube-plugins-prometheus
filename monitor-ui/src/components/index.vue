<template>
  <div class="index">
    <Menu></Menu>
    <div class="content">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import Menu from './sub-c/menu'
export default {
  name: 'index',
  methods: {
  },
  components:{
    Menu
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="less" scoped>
  @import '../assets/css/common.less';
  .index {
    min-width: 1200px;
    background: white;
  }
  .content {
    background-image: url('../assets/img/watermark_wecube_monitor.png');
    // width:1200px;
    min-height: 95vh;
    height: 100%;
    margin: 0 auto;
    padding: 10px;
    // background:white;
  }
</style>
