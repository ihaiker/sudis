<template>
    <e-charts class="w-100" :options="options"/>
</template>

<script>
    import ECharts from 'vue-echarts'
    import 'echarts/lib/chart/line'
    import utils from "../../tools/utils";

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
                    splitLine: {show: true},
                    nameGap: 0, boundaryGap: false,
                },
                yAxis: {
                    type: 'category', splitLine: {show: true},
                    axisLabel: {
                        formatter(limit, index) {
                            return utils.ram(limit);
                        }
                    },
                },
                series: [{
                    type: 'line', data: [],
                    showSymbol: false, hoverAnimation: false,
                }],
            },
        }),
        created() {
            this.options.series[0].data = this.lineData;
        },
    }
</script>
