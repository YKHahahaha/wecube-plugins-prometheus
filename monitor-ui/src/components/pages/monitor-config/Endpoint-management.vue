<template>
  <div class="main-content">
    <div v-if="showGroupMsg" style="padding-left:20px">
      <Tag type="border" closable color="primary" @on-close="closeTag">{{$t('field.group')}}:{{groupMsg.name}}</Tag>
    </div>
    <PageTable :pageConfig="pageConfig"></PageTable>
    <ModalComponent :modelConfig="modelConfig">
      <div slot="advancedConfig" class="extentClass">   
        <div class="marginbottom params-each">
          <label class="col-md-2 label-name lable-name-select">{{$t('field.endpoint')}}:</label>
          <Select v-model="modelConfig.slotConfig.resourceSelected" multiple filterable style="width:300px">
              <Option v-for="item in modelConfig.slotConfig.resourceOption" :value="item.id" :key="item.id">
                <Tag color="cyan" v-if="item.option_value.split(':')[1] == 'host'">host</Tag>
                <Tag color="blue" v-if="item.option_value.split(':')[1] == 'mysql'">mysql</Tag>
                <Tag color="geekblue" v-if="item.option_value.split(':')[1] == 'redis'">redis</Tag>
                <Tag color="purple" v-if="item.option_value.split(':')[1] == 'tomcat'">tomcat</Tag>
              {{ item.option_text }}</Option>
          </Select>
        </div>
      </div>
    </ModalComponent>
    <ModalComponent :modelConfig="historyAlarmModel">
      <div slot="historyAlarm">
        <tableTemp :table="historyAlarmPageConfig.table" :pageConfig="historyAlarmPageConfig"></tableTemp>
      </div>
    </ModalComponent>
    <ModalComponent :modelConfig="endpointRejectModel">
      <div slot="endpointReject">  
        <div class="marginbottom params-each">
          <label class="col-md-2 label-name lable-name-select">{{$t('field.endpoint')}}:</label>
          <Select v-model="endpointRejectModel.addRow.type" style="width:338px">
              <Option v-for="item in endpointRejectModel.endpointType" :value="item.value" :key="item.value">
              {{item.label}}</Option>
          </Select>
        </div>
        <div class="marginbottom params-each" v-if="showInstance">
          <label class="col-md-2 label-name lable-name-select">{{$t('field.instance')}}:</label>
          <input v-model="endpointRejectModel.addRow.instance" type="text" class="col-md-7 form-control model-input">
          <label class="required-tip">*</label>
        </div>
      </div>
    </ModalComponent>
  </div>
</template>
<script>
  import tableTemp from '@/components/components/table-page/table'
  let tableEle = [
    {title: 'tableKey.endpoint', value: 'guid', display: true},
    {title: 'tableKey.group', value: 'groups_name', display: true, }
  ]
  let historyAlarmEle = [
    {title: 'tableKey.status',value: 'status', style: 'min-width:70px', display: true},
    {title: 'tableKey.s_metric',value: 's_metric', display: true},
    {title: 'tableKey.start_value',value: 'start_value', display: true},
    {title: 'tableKey.s_cond',value: 's_cond', style: 'min-width:70px', display: true},
    {title: 'tableKey.s_last',value: 's_last', style: 'min-width:65px', display: true},
    {title: 'tableKey.s_priority',value: 's_priority', display: true},
    {title: 'tableKey.start',value: 'start_string', style: 'min-width:200px', display: true},
    {title: 'tableKey.end',value: 'end_string', style: 'min-width:200px',display: true,
      'render': (item) => {
        if (item.end === undefined) {
          return '-'
        } else {
          return item.end
        }
      }
    }]
  const btn = [
    {btn_name: 'button.thresholdManagement', btn_func: 'thresholdConfig'},
    {btn_name: 'button.historicalAlert', btn_func: 'historyAlarm'},
    {btn_name: 'button.remove', btn_func: 'delF'},
    {btn_name: 'button.logConfiguration', btn_func: 'logManagement'}
  ]
  export default {
    name: '',
    data() {
      return {
        pageConfig: {
          CRUD: this.apiCenter.endpointManagement.list.api,
          researchConfig: {
            input_conditions: [
              {value: 'search', type: 'input', placeholder: 'placeholder.input', style: ''}],
            btn_group: [
              {btn_name: 'button.search', btn_func: 'search', class: 'btn-confirm-f', btn_icon: 'fa fa-search'}
            ],
            filters: {
              search: ''
            }
          },
          table: {
            tableData: [],
            tableEle: tableEle,
            filterMoreBtn: 'filterMoreBtn',
            primaryKey: 'guid',
            btn: btn,
            pagination: this.pagination,
            handleFloat:true,
          },
          pagination: {
            __orders: '-created_date',
            total: 0,
            page: 1,
            size: 10
          }
        },
        historyAlarmPageConfig: {
          table: {
            tableData: [],
            tableEle: historyAlarmEle,
            btn: [],
          },
        },
        modelConfig: {
          modalId: 'add_object_Modal',
          modalTitle: 'button.historicalAlert',
          isAdd: true,
          config: [
            {name:'advancedConfig',type:'slot'}
          ],
          addRow: { // [通用]-保存用户新增、编辑时数据
            name: null,
            description: null,
          },
          slotConfig: {
            resourceSelected: [],
            resourceOption: []
          }
        },
        historyAlarmModel: {
          modalId: 'history_alarm_Modal',
          modalTitle: 'button.historicalAlert',
          modalStyle: 'width:930px;max-width: none;',
          noBtn: true,
          isAdd: true,
          config: [
            {name:'historyAlarm',type:'slot'}
          ],
          pageConfig: {
            table: {
              tableData: [],
              tableEle: tableEle
            }
          },
        },
        endpointRejectModel: {
          modalId: 'endpoint_reject_model',
          modalTitle: 'title.endpointAdd',
          isAdd: true,
          saveFunc: 'endpointRejectSave',
          config: [
            {name:'endpointReject',type:'slot'},
            {label: 'field.ip', value: 'exporter_ip', placeholder: 'tips.required', v_validate: 'required:true|isIP', disabled: false, type: 'text'},
            {label: 'field.port', value: 'exporter_port', placeholder: 'tips.required', v_validate: 'required:true|isNumber', disabled: false, type: 'text'},
          ],
          addRow: {
            instance: '',
            type: 'host',
            exporter_ip: null,
            exporter_port: 9100,
          },
          endpointType: [
            {label:'host',value:'host'},
            {label:'mysql',value:'mysql'},
            {label:'redis',value:'redis'},
            {label:'tomcat',value:'tomcat'}
          ],
        },
        id: null,
        showGroupMsg: false,
        groupMsg: {}
      }
    },
    mounted() {
      if (this.$validate.isEmpty_reset(this.$route.params)) {
        this.groupMsg = {}
        this.showGroupMsg = false
        this.pageConfig.researchConfig.btn_group.push({btn_name: 'button.add', btn_func: 'endpointReject', class: 'btn-cancle-f', btn_icon: 'fa fa-plus'})
      } else {
        this.$parent.activeTab = '/monitorConfigIndex/endpointManagement'
        if (this.$route.params.hasOwnProperty('group')) {
          this.groupMsg = this.$route.params.group
          this.showGroupMsg = true
          this.pageConfig.researchConfig.btn_group.push({btn_name: 'button.add', btn_func: 'add', class: 'btn-cancle-f', btn_icon: 'fa fa-plus'})
          this.pageConfig.researchConfig.filters.grp = this.groupMsg.id
        }
        if (this.$route.params.hasOwnProperty('search')) {
          this.pageConfig.researchConfig.filters.search = this.$route.params.search
        }
      }
      this.initData(this.pageConfig.CRUD, this.pageConfig)
    },
    watch: {
      'endpointRejectModel.addRow.type':function(val){
        const typeToPort = {
          host: 9100,
          mysql: 9104,
          redis: 9121,
          tomcat: 9151,
        }
        this.endpointRejectModel.addRow.exporter_port = typeToPort[val]
      }
    },
    computed:{
      showInstance: function(){
        return this.endpointRejectModel.addRow.type === 'host' ? false: true
      }
    },
    methods: {
      initData (url= this.pageConfig.CRUD, params) {
        this.$tableUtil.initTable(this, 'GET', url, params)
      },
      filterMoreBtn () {
        let moreBtnGroup = ['thresholdConfig','historyAlarm','logManagement']
        if (this.showGroupMsg) {
          moreBtnGroup.push('delF')
        }
        return moreBtnGroup
      },
      add () {
        this.modelConfig.slotConfig.resourceOption = []
        this.modelConfig.slotConfig.resourceSelected = []
        this.$httpRequestEntrance.httpRequestEntrance('GET',this.apiCenter.resourceSearch.api, {search: '.'}, responseData => {
          this.modelConfig.slotConfig.resourceOption = responseData
        })
        this.JQ('#add_object_Modal').modal('show')
      },
      addPost() {
        let params = {
          grp: this.groupMsg.id,
          endpoints: this.modelConfig.slotConfig.resourceSelected,
          operation: 'add'
        }
        this.$httpRequestEntrance.httpRequestEntrance('POST', this.apiCenter.endpointManagement.update.api, params, () => {
          this.$Message.success(this.$t('tips.success'))
          this.JQ('#add_object_Modal').modal('hide')
          this.initData(this.pageConfig.CRUD, this.pageConfig)
        })
      },
      delF (rowData) {
        let endpoints = []
        this.pageConfig.table.tableData.forEach((item)=>{
           endpoints.push(item.guid.split(':')[0])
        })
        let params = {
          grp: this.groupMsg.id,
          endpoints: [parseInt(rowData.id)],
          operation: 'delete'
        }
        this.$httpRequestEntrance.httpRequestEntrance('POST', this.apiCenter.endpointManagement.update.api, params, () => {
          this.$Message.success(this.$t('tips.success'))
          this.initData(this.pageConfig.CRUD, this.pageConfig)
        })
      },
      thresholdConfig (rowData) {
        this.$router.push({name: 'thresholdManagement', params: {id: rowData.id, type: 'endpoint'}})
      },
      logManagement (rowData) {
        this.$router.push({name: 'logManagement', params: {id: rowData.id, type: 'endpoint'}})
      },
      closeTag () {
        this.groupMsg = {}
        this.showGroupMsg = false
        this.pageConfig.researchConfig.filters.grp = ''
        this.pageConfig.table.btn.splice(this.pageConfig.table.btn.length-1, 1)
        this.pageConfig.researchConfig.btn_group.splice(this.pageConfig.researchConfig.btn_group.length-1, 1)
        this.pageConfig.researchConfig.btn_group.push({btn_name: 'button.add', btn_func: 'endpointReject', class: 'btn-cancle-f', btn_icon: 'fa fa-plus'})
        this.initData(this.pageConfig.CRUD, this.pageConfig)
      },
      historyAlarm (rowData) {
        let params = {id: rowData.id}
        this.$httpRequestEntrance.httpRequestEntrance('GET', this.apiCenter.alarm.history, params, (responseData) => {
          this.historyAlarmPageConfig.table.tableData = responseData
        })
        this.JQ('#history_alarm_Modal').modal('show')
      },
      endpointReject () {
        this.endpointRejectModel.addRow.type = 'host'
        this.JQ('#endpoint_reject_model').modal('show')
      },
      endpointRejectSave () {
        this.endpointRejectModel.addRow.exporter_port += ''
        let params= this.$validate.isEmptyReturn_JSON(this.endpointRejectModel.addRow)
        this.$httpRequestEntrance.httpRequestEntrance('POST', this.apiCenter.endpointManagement.register.api, params, () => {
          this.$validate.emptyJson(this.endpointRejectModel.addRow)
          this.JQ('#endpoint_reject_model').modal('hide')
          this.$Message.success(this.$t('tips.success'))
          this.initData(this.pageConfig.CRUD, this.pageConfig)
        })
      }
    },
    components: {
      tableTemp
    }
  }
</script>

<style lang="less" scoped>
</style>
