<template>
    <e-charts class="w-100" :options="options"/>
</template>

<script>
    import ECharts from 'vue-echarts'
    import 'echarts/lib/chart/line'

    export default {
        name: 'CPU', components: {ECharts},
        props: {
            lineData: {
                type: Array,
                default: [],
            }
        },
        data: () => ({
            options: {
                title: false, legend: false,
                tooltip : {
                    trigger: 'axis',
                },
                xAxis: {
                    type: 'category',
                    splitLine: {show: true}, nameGap: 0, boundaryGap: false,
                },
                yAxis: {
                    type: 'value',
                    axisLabel: {
                        formatter: '{value} %'
                    }, min: 0,
                },
                series: [{
                    type: 'line', data: [],
                    showSymbol: false, hoverAnimation: false,
                }],
            },
        }),
        mounted() {
            this.options.series[0].data = this.lineData;
        }
    }
</script>
