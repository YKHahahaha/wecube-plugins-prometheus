<template>
  <div class="pagination">
    <Page
      :total="pagination.total"
      :current="pagination.current"
      show-total show-sizer show-elevator
      placement="top"
      @on-change="pageChange"
      @on-page-size-change="pageSizeChange"
    ></Page>
  </div>
</template>
<script>
  export default {
    name: 'pagination',
    data () {
      return {
      }
    },
    props: ['pagination', 'pageUrl', 'filters'],
    mounted () {
    },
    methods: {
      pageChange (current) {
        // 换页时清空已选中数据
        this.$parent.clearSelectedData()
        this.pagination.page = current
        this.requestData()
      },
      pageSizeChange () {
        this.$parent.clearSelectedData()
        this.pagination.page = 1
        this.requestData()
      },
      requestData () {
        let pageConfig = {}
        pageConfig.researchConfig = {}
        pageConfig.researchConfig.filters = this.filters
        pageConfig.pagination = this.pagination
        this.$parent.initData(this.pageUrl, pageConfig)
      }
    },
    components: {}
  }
</script>

<style lang="less" scoped>
  .pagination{
    display: block;
    margin: 0px 0px 10px;
  }
</style>
