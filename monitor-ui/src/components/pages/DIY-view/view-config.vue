<template>
  <div class=" ">
      <header>
        <div style="display:flex;justify-content:space-between; font-size:16px;padding:8px 16px">
            <div class="header-name">
                <i class="fa fa-th-large fa-18" aria-hidden="true"></i>
                <span> {{$route.params.name}}</span>
            </div>
            <div class="header-tools"> 
              <button class="btn btn-sm btn-cancle-f" @click="addItem">{{$t('button.add')}}</button>
              <button class="btn btn-sm btn-confirm-f" @click="saveEdit">{{$t('button.saveEdit')}}</button>
              <button class="btn btn-sm btn-cancle-f" @click="goBack()">{{$t('button.back')}}</button>
            </div>
        </div>
      </header>
      <grid-layout
        :layout.sync="layoutData"
        :col-num="12"
        :row-height="30"
        :is-draggable="true"
        :is-resizable="true"
        :is-mirrored="false"
        :vertical-compact="true"
        :use-css-transforms="true"
        >
      <grid-item v-for="(item,index) in layoutData"
        :x="item.x"
        :y="item.y"
        :w="item.w"
        :h="item.h"
        :i="item.i"
        :key="index"
        @resize="resizeEvent"
        @resized="resizedEvent">
                 
        <div style="display:flex;justify-content:flex-end;padding:0 32px;">
          <div class="header-grid header-grid-name">
            <span v-if="editChartId !== item.id">{{item.i}}</span>
            <Input v-else v-model="item.i" class="editChartId" style="width:100px" @on-blur="editChartId = null" size="small" placeholder="small size" />
            <Tooltip :content="$t('placeholder.editTitle')" placement="top">
              <i class="fa fa-pencil-square" @click="editChartId = item.id" aria-hidden="true"></i>
            </Tooltip>
          </div>
          <div class="header-grid header-grid-tools"> 
            <!-- <i class="fa fa-eye" aria-hidden="true" @click="gridPlus(item)"></i> -->
            <Tooltip :content="$t('placeholder.chartConfiguration')" placement="top">
              <i class="fa fa-cog" @click="editGrid(item)" title="111" aria-hidden="true"></i>
            </Tooltip>
            <Tooltip :content="$t('placeholder.deleteChart')" placement="top">
              <i class="fa fa-trash" @click="removeGrid(item)" title="111" aria-hidden="true"></i>
            </Tooltip>
          </div>
        </div>
        <div class="">
          <section class="metric-section">
            <div v-if="!noDataTip">
              <div :id="item.id" class="echart" style="height:230px;width:560px"></div>
            </div>
             <div v-else class="echart echart-no-data-tip">
              <span>~~~No Data!~~~</span>
            </div>
          </section>
        </div>
        </grid-item>
      </grid-layout>
  </div>
</template>

<script>
import {generateUuid} from '@/assets/js/utils'
import VueGridLayout from 'vue-grid-layout'
import {drawChart} from  '@/assets/config/chart-rely'
const echarts = require('echarts/lib/echarts');
export default {
  name: '',
  data() {
    return {
      viewData: [],
      layoutData: [
        //   {'x':0,'y':0,'w':2,'h':2,'i':'0'},
        //   {'x':1,'y':1,'w':2,'h':2,'i':'1'},
      ],
      noDataTip: false,
      editChartId: null
    }
  },
  mounted() {
    if(this.$validate.isEmpty_reset(this.$route.params)) {
      this.$router.push({path:'viewConfigIndex'})
    } else {
      if (!this.$validate.isEmpty_reset(this.$route.params.cfg)) {
        this.viewData = JSON.parse(this.$route.params.cfg)
        this.initPanals()
      }
    }
  },
  methods: {
    initPanals () {
      this.viewData.forEach((item,viewIndex) => {
        this.layoutData.push(item.viewConfig)
        this.requestChart(item.viewConfig.id, item.query,viewIndex)
      })
    },
    requestChart (id, query,viewIndex) {
      let params = []
      query.forEach((item) => {
        params.push(JSON.stringify({
          endpoint: item.endpoint,
          metric: item.metricLabel,
          prom_ql: item.metric,
          time: '-1800'
        })) 
      })
      this.$httpRequestEntrance.httpRequestEntrance('GET',this.apiCenter.metricConfigView.api, {config: `[${params.join(',')}]`}, responseData => {
        var legend = []
        if (responseData.series.length === 0) {
          this.noDataTip = true
          return
        }
        const colorX = ['#66CC99','#9999CC','#9999FF','#99FF99','#CC9999','#CCCCFF','#FF66FF','#FFCCCC','#66CCCC','#996699','#99CC66','#FFCC99','#CC99CC','#FF9999','#FF9966','#FFCCFF']
        let colorSet = []
        for (let i=0;i<colorX.length;i++) {
          let tmpIndex = viewIndex*3 + i
          tmpIndex = tmpIndex%colorX.length
          colorSet.push(colorX[tmpIndex])
        }
        responseData.series.forEach((item,index)=>{
          legend.push(item.name)
          item.symbol = 'none'
          item.smooth = true
          item.lineStyle = {
            width: 1
          }
          item.itemStyle = {
            normal:{
              color: colorSet[index]
            }
          }
          item.areaStyle = {
            normal: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                  offset: 0,
                  color: colorSet[index]
              }, {
                  offset: 1,
                  color: 'white'
              }])
            }
          }
        }) 
        let config = {
          title: responseData.title,
          legend: legend,
          series: responseData.series,
          yaxis: responseData.yaxis
        }
        this.elId = id
        drawChart(this, config, {eye: false,dataZoom:false})
      })
    },
    addItem() {
      generateUuid().then((elId)=>{
        const key = ((new Date()).valueOf()).toString().substring(10)
        let item = {
          x:0,
          y:0,
          w:6,
          h:7,
          i: `default${key}`,
          id: `id_${elId}`
        }
        this.layoutData.push(item)
      })
    },
    editGrid(item) {
      this.modifyLayoutData().then((resViewData)=>{
        let parentRouteData = this.$route.params
        parentRouteData.cfg = JSON.stringify(resViewData) 
        this.$router.push({name: 'editView', params:{templateData: parentRouteData, panal:item, parentData: this.$route.params}}) 
      })
    },
    removeGrid(itemxxx) {
      this.layoutData.forEach((item,index) => {
        if (item.id === itemxxx.id) {
         this.layoutData.splice(index,1)
         return
        }
      })
    },
    gridPlus() {

    },
    modifyLayoutData() {
      return new Promise(resolve => {
        var resViewData = []
        this.layoutData.forEach((layoutDataItem) =>{
          let temp = {
            panalTitle: layoutDataItem.i,
            panalUnit: '',
            query: [],
            viewConfig: layoutDataItem
          }
          this.viewData.forEach((i) =>{
            if (layoutDataItem.id === i.viewConfig.id) {
              temp.panalUnit = i.panalUnit
              temp.query = i.query
            }
          })
          resViewData.push(temp)
        })
        resolve(resViewData)
      })
    },
    resizeEvent: function(i, newH, newW, newHPx, newWPx){
      this.layoutData.forEach((item,index) => {
        if (item.i === i) {
          this.layoutData[index].h = newH
          this.layoutData[index].w = newW
          var myChart = echarts.init(document.getElementById(item.id))
          myChart.resize({height:newHPx-64+'px',width:newWPx+'px'})
          return
        }
      })
    },
    resizedEvent: function(i, newH, newW, newHPx, newWPx){
      this.resizeEvent(i, newH, newW, newHPx, newWPx)
    },
    saveEdit() {
      let res = []
      this.layoutData.forEach((layoutDataItem) =>{
        this.viewData.forEach((i) =>{
          if (layoutDataItem.id === i.viewConfig.id) {
            res.push({
              panalTitle: i.panalTitle,
              panalUnit: i.panalUnit,
              query: i.query,
              viewConfig: layoutDataItem
            })
          }
        })
      })
      let params = {
        name: this.$route.params.name,
        id: this.$route.params.id,
        cfg: JSON.stringify(res)
      }
      this.$httpRequestEntrance.httpRequestEntrance('POST',this.apiCenter.template.save, params, () => {
        this.$Message.success(this.$t('tips.success'))
      })
    },
    goBack () {
      this.$router.push({name:'viewConfigIndex'})
    },
  },
  components: {
    GridLayout: VueGridLayout.GridLayout,
    GridItem: VueGridLayout.GridItem,
  },
}
</script>

<style scoped lang="less">
  .header-grid {
    flex-grow: 1;
    text-align: end;
    line-height: 32px;
    i {
      margin: 0 4px;
      cursor: pointer;
    } 
  }
</style>
<style scoped lang="less">

.vue-grid-item:not(.vue-grid-placeholder) {
    background: @gray-f;
    border: 1px solid @gray-f;
}
.echart-no-data-tip {
  text-align: center;
}
</style>
