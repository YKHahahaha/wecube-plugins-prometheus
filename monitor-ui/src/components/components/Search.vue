<template>
  <div class="" style="display:inline-block">
   <ul class="search-ul">
      <li class="search-li">
        <Select
          style="width:300px;"
          v-model="endpoint"
          filterable
          clearable
          remote
          :placeholder="$t('placeholder.input')"
          :remote-method="getEndpointList"
          @on-clear="clearEndpoint"
          >
          <Option v-for="(option, index) in endpointList" :value="option.option_value" :key="index">
            <Tag color="cyan" class="tag-width" v-if="option.option_value.split(':')[1] == 'host'">host</Tag>
            <Tag color="blue" class="tag-width" v-if="option.option_value.split(':')[1] == 'mysql'">mysql </Tag>
            <Tag color="geekblue" class="tag-width" v-if="option.option_value.split(':')[1] == 'redis'">redis </Tag>
            <Tag color="purple" class="tag-width" v-if="option.option_value.split(':')[1] == 'tomcat'">tomcat</Tag>{{option.option_text}}</Option>
        </Select>
      </li>
      <li class="search-li">
        <button type="button" class="btn btn-sm btn-confirm-f"
            @click="getChartsConfig()">
            <i class="fa fa-search" ></i>
            {{$t('button.search')}}
          </button>
      </li>
      <li class="search-li">
          <Select v-model="timeTnterval" style="width:80px" @on-change="getChartsConfig">
            <Option v-for="item in dataPick" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
      </li>
      <li class="search-li">
        <DatePicker type="daterange" placement="bottom-end" @on-change="datePick" :placeholder="$t('placeholder.datePicker')" style="width: 200px"></DatePicker>
      </li>
      <li class="search-li">
        <Select v-model="autoRefresh" style="width:100px" @on-change="getChartsConfig" :placeholder="$t('placeholder.refresh')">
          <Option v-for="item in autoRefreshConfig" :value="item.value" :key="item.value">{{ item.label }}</Option>
        </Select>
      </li>
      <li class="search-li">
        <button type="button" v-if="isShow" @click="changeRoute" class="btn btn-sm btn-cancle-f btn-jump">{{$t('button.endpointManagement')}}</button>
      </li>
   </ul>
  </div>
</template>

<script>
import {dataPick, autoRefreshConfig} from '@/assets/config/common-config'
export default {
  name: '',
  data() {
    return {
      endpoint: '',
      endpointList: [],
      searchInputConfig: {
        poptipWidth: 300,
        placeholder: 'placeholder.endpointSearch',
        inputStyle: "width:300px;",
        api: this.apiCenter.resourceSearch.api
      },
      ip: {},
      timeTnterval: -1800,
      dataPick: dataPick,
      dateRange: ['',''],
      autoRefresh: 0,
      autoRefreshConfig: autoRefreshConfig,
      params: {
        // time: this.timeTnterval,
        // group: 1,
        // endpoint: '192.168.0.16',
        // start: Date.parse(this.dateRange[0]),
        // end: Date.parse(this.dateRange[1])
      }
    }
  },
  computed: {
    isShow: function () {
      return !this.$validate.isEmpty_reset(this.endpoint)
    }
  },
  mounted() {
    if (!this.$validate.isEmpty_reset(this.$route.params)) {
      this.getEndpointList('.')
      this.endpoint = this.$route.params.value
    }
  },
  methods: {
    getMainConfig () {
      return new Promise(resolve => {
        let params = {
          type: this.endpoint.split(':')[1]
        }
        this.$httpRequestEntrance.httpRequestEntrance('GET', this.apiCenter.mainConfig.api, params, (responseData) => {
          resolve(responseData)
        })
      })
    },
    datePick (data) {
      this.dateRange = data
      if (this.dateRange[0] !== '') {
        this.dateRange[0] = this.dateRange[0] + ' 00:00:01'
      }
      if (this.dateRange[1] !== '') {
        this.dateRange[1] = this.dateRange[1] + ' 23:59:59'
      }
      this.getChartsConfig()
    },
    getEndpointList(query) {
      let params = {
        search: query,
        page: 1,
        size: 1000
      }
      this.$httpRequestEntrance.httpRequestEntrance('GET', this.apiCenter.resourceSearch.api, params, (responseData) => {
       this.endpointList = responseData
      })
    },
    getChartsConfig () {
      if (this.$validate.isEmpty_reset(this.endpoint)) {
        return
      }
      this.getMainConfig().then((res)=>{
        let url = res.panels.url
        let key = res.search.name
        let params = {
          autoRefresh: this.autoRefresh,
          time: this.timeTnterval,
          endpoint: this.endpoint.split(':')[0],
          start: this.dateRange[0] ===''? '':Date.parse(this.dateRange[0])/1000,
          end: this.dateRange[1] ===''? '':Date.parse(this.dateRange[1])/1000
        }
        url = url.replace(`{${key}}`,params[key].split(':')[0])
        this.$httpRequestEntrance.httpRequestEntrance('GET',url, params, responseData => {
          this.$parent.manageCharts(responseData, params)
        },{isNeedloading: false})
      })
    },
    clearEndpoint () {
      this.clearEndpoint = []
      this.$parent.showCharts = false
    },
    changeRoute () {
      this.$router.push({name: 'endpointManagement', params: {search: this.endpoint.split(':')[0]}})
    }
  },
  components: {
  }
}
</script>

<style scoped lang="less">
  .search-li {
    display: inline-block;
  }
  .search-ul>li:not(:first-child) {
    padding-left: 10px;
  }
  .tag-width {
    width: 55px;
    text-align: center;
  }
</style>
