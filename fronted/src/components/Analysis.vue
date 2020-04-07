<template>
  <div class="content-wrapper">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0 text-dark">统计分析</h1>
        </div>
        <div class="col-sm-6">
          <div class="form-group">
            <div class="input-group">
              <div class="input-group-prepend">
                      <span class="input-group-text">
                        <i class="far fa-calendar-alt"></i>
                      </span>
              </div>
              <input type="text" class="form-control float-right" id="reservation">
            </div>
          </div>
        </div>
      </div>
    </div>
    <section class="content">
      <div class="container-fluid">
        <div class="row">
          <div class="col-lg-6 col-9">
            <div class="small-box bg-info">
              <div class="inner">
                <h3 id="allocate_total">{{analysis_allocate_total}}</h3>
                <p>总分配人数</p>
              </div>
              <div class="icon">
                <i class="ion ion-person-add"></i>
              </div>
            </div>
          </div>
          <div class="col-lg-6 col-9">
            <div class="small-box bg-success">
              <div class="inner">
                <h3 id="activity_total">{{analysis_activity_total}}</h3>
                <p>总活跃用户数</p>
              </div>
              <div class="icon">
                <i class="ion ion-stats-bars"></i>
              </div>
            </div>
          </div>

          <div class="col-md-6">
            <div class="card card-danger">
              <div class="card-header">
                <h3 class="card-title">分组统计图示1</h3>
                <div class="card-tools">
                  <button type="button" class="btn btn-tool" data-card-widget="collapse"><i class="fas fa-minus"></i>
                  </button>
                  <button type="button" class="btn btn-tool" data-card-widget="remove"><i class="fas fa-times"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div id="group_pie"></div>
              </div>
            </div>
            <!-- /.card -->

            <!-- PIE CHART -->
            <div class="card card-danger">
              <div class="card-header">
                <h3 class="card-title">分域统计图示1</h3>

                <div class="card-tools">
                  <button type="button" class="btn btn-tool" data-card-widget="collapse"><i class="fas fa-minus"></i>
                  </button>
                  <button type="button" class="btn btn-tool" data-card-widget="remove"><i class="fas fa-times"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div id="domain_bar"></div>
              </div>
              <!-- /.card-body -->
            </div>
            <!-- /.card -->

          </div>
          <!-- /.col (LEFT) -->
          <div class="col-md-6">


            <!-- BAR CHART -->
            <div class="card card-success">
              <div class="card-header">
                <h3 class="card-title">分组统计图示2</h3>

                <div class="card-tools">
                  <button type="button" class="btn btn-tool" data-card-widget="collapse"><i class="fas fa-minus"></i>
                  </button>
                  <button type="button" class="btn btn-tool" data-card-widget="remove"><i class="fas fa-times"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div class="chart">
                  <!--                  <canvas id="group_bar"-->
                  <!--                          style="min-height: 250px; height: 250px; max-height: 250px; max-width: 100%;"></canvas>-->
                                    <div id="group_bar"></div>

<!--                  <x-chart :id="id" :option="option"></x-chart>-->
                </div>
              </div>
              <!-- /.card-body -->
            </div>
            <!-- /.card -->

            <!-- STACKED BAR CHART -->
            <div class="card card-success">
              <div class="card-header">
                <h3 class="card-title">分域统计图示2</h3>

                <div class="card-tools">
                  <button type="button" class="btn btn-tool" data-card-widget="collapse"><i class="fas fa-minus"></i>
                  </button>
                  <button type="button" class="btn btn-tool" data-card-widget="remove"><i class="fas fa-times"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div class="chart">
                  <div id="domain_pie"></div>
                </div>
              </div>
              <!-- /.card-body -->
            </div>
            <!-- /.card -->

          </div>
          <!-- /.col (RIGHT) -->
        </div>
        <!-- /.row -->
      </div><!-- /.container-fluid -->
    </section>
    <!-- /.content -->
  </div>
  <!-- /.content-wrapper -->
</template>

<script>
  import Highcharts from 'highcharts/highstock';
  import HighchartsMore from 'highcharts/highcharts-more';
  import HighchartsDrilldown from 'highcharts/modules/drilldown';
  import Highcharts3D from 'highcharts/highcharts-3d';
  import Highmaps from 'highcharts/modules/map';
  import $ from 'jquery'

  HighchartsMore(Highcharts)
  HighchartsDrilldown(Highcharts);
  Highcharts3D(Highcharts);
  Highmaps(Highcharts);

  import XChart from './chart.vue'

  export default {
    name: 'HelloWorld',

    props: {
      msg: String
    },
    data() {
      return{
        id: 'test',
        analysis_allocate_total:0,
        analysis_activity_total:0
      }
    },
    components: {
      XChart
    },
    methods:{},

    mounted() {
      this.$axios.get('/apis/analysis_allocate_total')
        .then(response => {
          this.analysis_allocate_total = response.data[0]['count'];
        })
        .catch(error => {
          console.log(error);
          alert('网络错误，不能访问');
        })

      this.$axios.get('/apis/analysis_activity_total')
        .then(response => {
          this.analysis_activity_total = response.data[0]['total'];
        })
        .catch(error => {
          console.log(error);
          alert('网络错误，不能访问');
        })


      this.$axios.get('/apis/analysis_auth_group')
        .then(response => {
          var total=0
          var series=[]

          var labels=[]
          var counts=[]
          for (var i = 0; i < response.data.length; i++) {
            total=total+(response.data)[i]['total']
            labels.push((response.data)[i]['group_id'])
            counts.push((response.data)[i]['total'])
          }
          for (var i = 0; i < response.data.length; i++) {

            series.push({name:(response.data)[i]['group_id'],y:(response.data)[i]['total']/total})

          }


         var group_pie= {
            chart: {
              plotBackgroundColor: null,
              plotBorderWidth: null,
              plotShadow: false,
              type: 'pie'
            },
            title: {
              text: '分组授权所占比例'
            },
            tooltip: {
              pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
            },
            plotOptions: {
              pie: {
                allowPointSelect: true,
                cursor: 'pointer',
                dataLabels: {
                  enabled: true,
                  format: '<b>{point.name}</b>: {point.percentage:.1f} %',
                  style: {
                    // color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                  }
                }
              }
            },
            series: [{
              name: 'Brands',
              colorByPoint: true,
              data: series
            }]
          }
          new Highcharts.chart("group_bar",group_pie)




          var group_bar= {
            chart: {
              type: 'column'
            },
            title: {
              text: '分组授权总数'
            },
            xAxis: {
              categories: labels,
              crosshair: true
            },
            yAxis: {
              min: 0,
              title: {
                text: '授权次数'
              }
            },
            tooltip: {
              // head + 每个 point + footer 拼接成完整的 table
              headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
              pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
                '<td style="padding:0"><b>{point.y:.1f} mm</b></td></tr>',
              footerFormat: '</table>',
              shared: true,
              useHTML: true
            },
            plotOptions: {
              column: {
                borderWidth: 0
              }
            },
            series: [{
              name: '访问量',
              data: counts
            }]
          }
          new Highcharts.chart("group_pie",group_bar)

        })
        .catch(error => {
          console.log(error);
          alert('网络错误，不能访问');
        })

      this.$axios.get('/apis/analysis_use_group')
        .then(response => {
          var total=0
          var series=[]

          var labels=[]
          var counts=[]
          for (var i = 0; i < response.data.length; i++) {
            total=total+(response.data)[i]['cgulaid']
            labels.push((response.data)[i]['group_name'])
            counts.push((response.data)[i]['cgulaid'])
          }
          for (var i = 0; i < response.data.length; i++) {

            series.push({name:(response.data)[i]['group_name'],y:(response.data)[i]['cgulaid']/total})

          }


          var group_pie= {
            chart: {
              plotBackgroundColor: null,
              plotBorderWidth: null,
              plotShadow: false,
              type: 'pie'
            },
            title: {
              text: '2018年1月浏览器市场份额'
            },
            tooltip: {
              pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
            },
            plotOptions: {
              pie: {
                allowPointSelect: true,
                cursor: 'pointer',
                dataLabels: {
                  enabled: true,
                  format: '<b>{point.name}</b>: {point.percentage:.1f} %',
                  style: {
                    // color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                  }
                }
              }
            },
            series: [{
              name: 'Brands',
              colorByPoint: true,
              data: series
            }]
          }
          new Highcharts.chart("domain_bar",group_pie)




          var group_bar= {
            chart: {
              type: 'column'
            },
            title: {
              text: '月平均降雨量'
            },
            subtitle: {
              text: '数据来源: WorldClimate.com'
            },
            xAxis: {
              categories: labels,
              crosshair: true
            },
            yAxis: {
              min: 0,
              title: {
                text: '降雨量 (mm)'
              }
            },
            tooltip: {
              // head + 每个 point + footer 拼接成完整的 table
              headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
              pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
                '<td style="padding:0"><b>{point.y:.1f} mm</b></td></tr>',
              footerFormat: '</table>',
              shared: true,
              useHTML: true
            },
            plotOptions: {
              column: {
                borderWidth: 0
              }
            },
            series: [{
              name: '访问量',
              data: counts
            }]
          }
          new Highcharts.chart("domain_pie",group_bar)

        })
        .catch(error => {
          console.log(error);
          alert('网络错误，不能访问');
        })

      // var domain_pie_data = {
      //   labels:  this.analysis_use_domain_lables,
      //   datasets: [
      //     {
      //       data: this.analysis_use_domain_datas,
      //       backgroundColor: ['#f56954', '#00a65a', '#f39c12'],
      //     }
      //   ]
      // }
      // this.pie_data.datasets.data=eval(this.analysis_use_domain_datas)
      // this.pie_data.labels=eval(this.analysis_use_domain_lables)

      // console.log(domain_pie_obj)
      // console.log(eval((this.pie_data)))
      // new Chart(domain_pie, {
      //   type: 'doughnut',
      //   data: this.pie_data,
      //   options: {
      //     maintainAspectRatio: false,
      //     responsive: true,
      //   }
      // })


      //
      // var domain_bar_data = {
      //   labels: this.analysis_use_domain_lables,
      //   datasets: [
      //     {
      //       label: '域名',
      //       backgroundColor: 'rgba(60,141,188,0.9)',
      //       borderColor: 'rgba(60,141,188,0.8)',
      //       pointRadius: false,
      //       pointColor: '#3b8bba',
      //       pointStrokeColor: 'rgba(60,141,188,1)',
      //       pointHighlightFill: '#fff',
      //       pointHighlightStroke: 'rgba(60,141,188,1)',
      //       data: this.analysis_use_domain_datas
      //     }
      //   ]
      // }
      //
      // var domain_bar_option = {
      //   responsive: true,
      //   maintainAspectRatio: false,
      //   datasetFill: false
      // }
      //
      // new Chart(doamin_bar_obj, {
      //   type: 'bar',
      //   data: domain_bar_data,
      //   options: domain_bar_option
      // })

      //
      // var group_pie_obj = $('#group_pie').get(0).getContext('2d')
      // var group_pie_data = {
      //   labels: this.analysis_use_group_labels,
      //   datasets: [
      //     {
      //       data: this.analysis_use_group_datas,
      //       backgroundColor: ['#f56954', '#00a65a', '#f39c12', '#00c0ef', '#3c8dbc', '#d2d6de', '#d2d6df', '#d2d6dd'],
      //     }
      //   ]
      // }
      //
      // var group_pie_option = {
      //   maintainAspectRatio: false,
      //   responsive: true,
      // }
      // new Chart(group_pie_obj, {
      //   type: 'doughnut',
      //   data: group_pie_data,
      //   options: group_pie_option
      // })
      //
      //
      //
      //
      // var group_bar_data = {
      //   labels: this.analysis_use_group_labels,
      //   datasets: [
      //     {
      //       label: '组名',
      //       backgroundColor: 'rgba(60,141,188,0.9)',
      //       borderColor: 'rgba(60,141,188,0.8)',
      //       pointRadius: false,
      //       pointColor: '#3b8bba',
      //       pointStrokeColor: 'rgba(60,141,188,1)',
      //       pointHighlightFill: '#fff',
      //       pointHighlightStroke: 'rgba(60,141,188,1)',
      //       data: this.analysis_use_group_datas
      //     }
      //   ]
      // }
      // var group_bar_obj = $('#group_bar').get(0).getContext('2d')
      // var group_bar_option = {
      //   responsive: true,
      //   maintainAspectRatio: false,
      //   datasetFill: false
      // }
      //
      // new Chart(group_bar_obj, {
      //   type: 'bar',
      //   data: group_bar_data,
      //   options: group_bar_option
      // })
    },
    error: function (e) {
      console.log(e.status);
      console.log(e.responseText);
    }


  }

</script>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h3 {
    margin: 40px 0 0;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }
</style>
